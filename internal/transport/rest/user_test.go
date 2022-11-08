package rest

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	mocks_handler "github.com/Arkosh744/InternAvitoBackend/internal/transport/rest/mocks"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/types"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/validator"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Create(t *testing.T) {
	// Init Test Table

	type mockBehavior func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.InputUser, out domain.User)
	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domain.InputUser
		out                  domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"firstName": "TestName", "email": "test@test.test", "lastName": "TestLastName"}`,
			inputUser: domain.InputUser{
				FirstName: "TestName",
				LastName:  "TestLastName",
				Email:     "test@test.test",
			},
			out: domain.User{
				FirstName: "TestName",
				LastName:  "TestLastName",
				Email:     "test@test.test",
				ID:        uuid.MustParse("7f6280a8-64fb-4408-b16b-1e77b9fa5f35"),
			},
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.InputUser, out domain.User) {
				r.EXPECT().CheckUserByEmail(gomock.Any(), inp.Email).Return(domain.User{}, nil)
				r.EXPECT().Create(gomock.Any(), inp).Return(out, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":"7f6280a8-64fb-4408-b16b-1e77b9fa5f35","firstName":"TestName","lastName":"TestLastName","email":"test@test.test"}`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"firstName": "TestName", "lastName": "TestLastName"}`,
			inputUser: domain.InputUser{
				FirstName: "TestName",
				LastName:  "TestLastName",
			},
			out: domain.User{
				FirstName: "TestName",
				LastName:  "TestLastName",
				Email:     "tesst.test",
				ID:        uuid.MustParse("7f6280a8-64fb-4408-b16b-1e77b9fa5f35"),
			},
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.InputUser, out domain.User) {

			},
			expectedStatusCode:   422,
			expectedResponseBody: `{"message":"Key: 'InputUser.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{
			name:      "Already Exists",
			inputBody: `{"firstName": "TestName", "email": "test@test.ru", "lastName": "TestLastName"}`,
			inputUser: domain.InputUser{
				FirstName: "TestName",
				LastName:  "TestLastName",
				Email:     "test@test.ru",
			},
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.InputUser, out domain.User) {
				r.EXPECT().CheckUserByEmail(gomock.Any(), inp.Email).Return(domain.User{}, types.ErrAlreadyExists)
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"message":"could not create user: already exists"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"firstName": "TestName", "email": "test@test.ru", "lastName": "TestLastName"}`,
			inputUser: domain.InputUser{
				FirstName: "TestName",
				LastName:  "TestLastName",
				Email:     "test@test.ru",
			},
			out: domain.User{},
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.InputUser, out domain.User) {
				r.EXPECT().CheckUserByEmail(gomock.Any(), inp.Email).Return(domain.User{}, nil)
				r.EXPECT().Create(gomock.Any(), inp).Return(domain.User{}, types.ErrorInternal)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"could not create user: internal error"}`,
		},
	}

	// Run Tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mocks_handler.NewMockUsers(c)
			test.mockBehavior(users, context.Background(), test.inputUser, test.out)
			h := Handler{usersService: users}
			// Init Endpoint
			router := echo.New()
			router.Validator = validator.NewValidator()

			v1 := router.Group("/v1")
			userRoutes := v1.Group("/user")
			userRoutes.POST("/", h.Create)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/v1/user/", bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			// Make Request
			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, strings.Split(w.Body.String(), "\n")[0], test.expectedResponseBody)

		})
	}
}

func TestHandler_GetUserBalance(t *testing.T) {
	// Init Test Table

	type mockBehavior func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.User, out domain.User)
	userID := "7f6280a8-64fb-4408-b16b-1e77b9fa5f35"
	tests := []struct {
		name                 string
		inputID              string
		out                  domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			inputID: userID,
			out: domain.User{
				FirstName: "TestName",
				LastName:  "TestLastName",
				Email:     "test@test.test",
				Wallet:    domain.Wallet{Balance: 500},
				ID:        uuid.MustParse(userID),
			},
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.User, out domain.User) {
				r.EXPECT().GetUserBalance(gomock.Any(), inp).Return(inp, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message"": Balance of test@test.test, "balance": "500"}`,
		},
		{
			name:    "Wrong Input",
			inputID: userID,
			out: domain.User{
				FirstName: "TestName",
				LastName:  "TestLastName",
				Email:     "test@test.test",
				ID:        uuid.MustParse(userID),
			},
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.User, out domain.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid UUID length: 5"}`,
		},
		{
			name:    "Service Error",
			inputID: userID,
			out:     domain.User{},
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp domain.User, out domain.User) {
				r.EXPECT().GetUserBalance(gomock.Any(), inp).Return(inp, types.ErrorInternal)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"could not get user balance internal error"}`,
		},
	}

	// Run Tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mocks_handler.NewMockUsers(c)
			test.mockBehavior(users, context.Background(), domain.User{ID: uuid.MustParse(test.inputID)}, test.out)
			h := Handler{usersService: users}
			// Init Endpoint
			router := echo.New()
			router.Validator = validator.NewValidator()

			reqID := fmt.Sprintf("/v1/user/%v", test.inputID)
			if test.name == "Wrong Input" {
				reqID = fmt.Sprintf("/v1/user/%v", "wrong")
			}
			v1 := router.Group("/v1")
			userRoutes := v1.Group("/user")
			userRoutes.GET("/:id", h.GetUserBalance)

			// Create Request
			w := httptest.NewRecorder()
			fmt.Println(reqID)
			req := httptest.NewRequest("GET", reqID, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			// Make Request
			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			if test.name == "Ok" {
				assert.Equal(t, fmt.Sprintf(`{"message"": Balance of %s, "balance": "%v"}`,
					test.out.Email, test.out.Wallet.Balance), test.expectedResponseBody)
			} else {
				assert.Equal(t, strings.Split(w.Body.String(), "\n")[0], test.expectedResponseBody)
			}
		})
	}
}
