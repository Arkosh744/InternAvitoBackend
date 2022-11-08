package rest

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
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

func TestHandler_BuyServiceUser(t *testing.T) {
	// Init Test Table

	type mockBehavior func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputBuyServiceUser)
	IDUser := "7f6280a8-64fb-4408-b16b-1e77b9fa5f35"

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            wallet.InputBuyServiceUser
		balance              float32
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","cost": 100,"service_name": "Google Eda"}`,
			inputUser: wallet.InputBuyServiceUser{
				IDUser:      uuid.MustParse(IDUser),
				ServiceName: "Google Eda",
				Cost:        100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputBuyServiceUser) {
				r.EXPECT().BuyServiceUser(gomock.Any(), inp).Return(
					wallet.OutPendingOrder{ID: uuid.MustParse("2d0cefd9-034d-4af9-a86c-70ac8efd4842"),
						Cost: inp.Cost, ServiceName: inp.ServiceName, Status: "created",
						Txn: uuid.MustParse("7fe842c5-94c8-4bb0-92bd-3f881cac4f34")}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Created order for Google Eda for 100","order_id":"2d0cefd9-034d-4af9-a86c-70ac8efd4842"}`,
		}, {
			name:      "Internal_error",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","cost": 100,"service_name": "Google Eda"}`,
			inputUser: wallet.InputBuyServiceUser{
				IDUser:      uuid.MustParse(IDUser),
				ServiceName: "Google Eda",
				Cost:        100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputBuyServiceUser) {
				r.EXPECT().BuyServiceUser(gomock.Any(), inp).Return(
					wallet.OutPendingOrder{}, types.ErrorInternal)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"internal error"}`,
		},
	}

	// Run Tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mocks_handler.NewMockUsers(c)
			test.mockBehavior(users, context.Background(), test.inputUser)
			h := Handler{usersService: users}
			router := echo.New()
			router.Validator = validator.NewValidator()

			v1 := router.Group("/v1")
			userRoutes := v1.Group("/user")
			userRoutes.POST("/", h.Create)
			// Init Endpoint
			walletRoutes := userRoutes.Group("/wallet")
			orderRoutes := walletRoutes.Group("/order")
			orderRoutes.POST("/buy", h.BuyServiceUser)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/v1/user/wallet/order/buy", bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			// Make Request
			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, strings.Split(w.Body.String(), "\n")[0], test.expectedResponseBody)
		})
	}
}

func TestHandler_ManageOrder(t *testing.T) {
	// Init Test Table

	type mockBehavior func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputOrderManager)
	IDUser := "7f6280a8-64fb-4408-b16b-1e77b9fa5f35"
	IDOrder := "2d0cefd9-034d-4af9-a86c-70ac8efd4842"

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            wallet.InputOrderManager
		balance              float32
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Approve_ok",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","id_order": "2d0cefd9-034d-4af9-a86c-70ac8efd4842"}`,
			inputUser: wallet.InputOrderManager{
				IDUser:  uuid.MustParse(IDUser),
				IDOrder: uuid.MustParse(IDOrder),
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputOrderManager) {
				r.EXPECT().ManageOrder(gomock.Any(), wallet.InputOrderManager{IDUser: inp.IDUser, IDOrder: inp.IDOrder, Status: "approved"}).Return(
					wallet.OutOrderManager{IDOrder: uuid.MustParse(IDOrder), Status: "completed"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Order completed"}`,
		}, {
			name:      "Decline_ok",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","id_order": "2d0cefd9-034d-4af9-a86c-70ac8efd4842"}`,
			inputUser: wallet.InputOrderManager{
				IDUser:  uuid.MustParse(IDUser),
				IDOrder: uuid.MustParse(IDOrder),
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputOrderManager) {
				r.EXPECT().ManageOrder(gomock.Any(), wallet.InputOrderManager{IDUser: inp.IDUser, IDOrder: inp.IDOrder, Status: "cancelled"}).Return(
					wallet.OutOrderManager{IDOrder: uuid.MustParse(IDOrder), Status: "cancelled"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Order cancelled"}`,
		}, {
			name:      "Decline_internal_error",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","id_order": "2d0cefd9-034d-4af9-a86c-70ac8efd4842"}`,
			inputUser: wallet.InputOrderManager{
				IDUser:  uuid.MustParse(IDUser),
				IDOrder: uuid.MustParse(IDOrder),
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputOrderManager) {
				r.EXPECT().ManageOrder(gomock.Any(),
					wallet.InputOrderManager{IDUser: inp.IDUser, IDOrder: inp.IDOrder, Status: "cancelled"}).
					Return(wallet.OutOrderManager{}, types.ErrorInternal)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"internal error"}`,
		},
	}

	// Run Tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mocks_handler.NewMockUsers(c)
			test.mockBehavior(users, context.Background(), test.inputUser)
			h := Handler{usersService: users}
			router := echo.New()
			router.Validator = validator.NewValidator()

			v1 := router.Group("/v1")
			userRoutes := v1.Group("/user")
			userRoutes.POST("/", h.Create)
			// Init Endpoint
			walletRoutes := userRoutes.Group("/wallet")
			orderRoutes := walletRoutes.Group("/order")
			orderRoutes.POST("/approve", h.ManageOrder)
			orderRoutes.POST("/decline", h.ManageOrder)
			var path string
			if strings.Split(test.name, "_")[0] == "Approve" {
				path = fmt.Sprintf("/v1/user/wallet/order/approve")
			} else if strings.Split(test.name, "_")[0] == "Decline" {
				path = fmt.Sprintf("/v1/user/wallet/order/decline")
			}

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", path, bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			// Make Request
			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, strings.Split(w.Body.String(), "\n")[0], test.expectedResponseBody)
		})
	}
}
