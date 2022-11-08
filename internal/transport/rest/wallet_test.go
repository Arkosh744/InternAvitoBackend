package rest

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
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

func TestHandler_DepositToUser(t *testing.T) {
	// Init Test Table

	type mockBehavior func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit)
	IDUser := "7f6280a8-64fb-4408-b16b-1e77b9fa5f35"
	IDWallet := "3ad280a8-64fb-4128-b16b-1e77b9f33f35"

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            wallet.InputDeposit
		balance              float32
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ID_user_no_wallet",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","amount": 100}`,
			inputUser: wallet.InputDeposit{
				IDUser:   uuid.MustParse(IDUser),
				IDWallet: uuid.Nil,
				Amount:   100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByUserID(gomock.Any(), inp.IDUser).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{ID: uuid.Nil}}, nil)
				r.EXPECT().CreateWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet)}}, nil)
				inp.IDWallet = uuid.MustParse(IDWallet)
				r.EXPECT().DepositWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet), Balance: 100}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Deposited 100 to test@tgest.ru. Balance: 100"}`,
		}, {
			name:      "ID_user_with_wallet",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","amount": 100}`,
			inputUser: wallet.InputDeposit{
				IDUser:   uuid.MustParse(IDUser),
				IDWallet: uuid.Nil,
				Amount:   100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByUserID(gomock.Any(), inp.IDUser).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet)}}, nil)
				inp.IDWallet = uuid.MustParse(IDWallet)
				r.EXPECT().DepositWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet), Balance: 100}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Deposited 100 to test@tgest.ru. Balance: 100"}`,
		}, {
			name:      "email_user_with_wallet",
			inputBody: `{"email": "test@tgest.ru","amount": 100}`,
			inputUser: wallet.InputDeposit{
				EmailUser: "test@tgest.ru",
				IDUser:    uuid.Nil,
				IDWallet:  uuid.Nil,
				Amount:    100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByEmail(gomock.Any(), inp.EmailUser).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet)}}, nil)
				inp.IDWallet = uuid.MustParse(IDWallet)
				r.EXPECT().DepositWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet), Balance: 100}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Deposited 100 to test@tgest.ru. Balance: 100"}`,
		}, {
			name:      "email_user_no_wallet",
			inputBody: `{"email": "test@tgest.ru","amount": 100}`,
			inputUser: wallet.InputDeposit{
				EmailUser: "test@tgest.ru",
				IDUser:    uuid.Nil,
				IDWallet:  uuid.Nil,
				Amount:    100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByEmail(gomock.Any(), inp.EmailUser).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.Nil}}, nil)
				r.EXPECT().CreateWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet)}}, nil)
				inp.IDWallet = uuid.MustParse(IDWallet)
				r.EXPECT().DepositWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet), Balance: 100}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Deposited 100 to test@tgest.ru. Balance: 100"}`,
		}, {
			name:      "ID_wallet",
			inputBody: `{"id_wallet": "3ad280a8-64fb-4128-b16b-1e77b9f33f35","amount": 100}`,
			inputUser: wallet.InputDeposit{
				IDUser:   uuid.Nil,
				IDWallet: uuid.MustParse(IDWallet),
				Amount:   100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().DepositWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet), Balance: 100}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Deposited 100 to test@tgest.ru. Balance: 100"}`,
		}, {
			name:      "wrong_email",
			inputBody: `{"email": "test","amount": 100}`,
			inputUser: wallet.InputDeposit{
				EmailUser: "test",
				Amount:    100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
			},
			expectedStatusCode:   422,
			expectedResponseBody: `{"message":"Key: 'InputDeposit.EmailUser' Error:Field validation for 'EmailUser' failed on the 'email' tag"}`,
		}, {
			name:      "No_user",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","amount": 100}`,
			inputUser: wallet.InputDeposit{
				IDUser:   uuid.MustParse(IDUser),
				IDWallet: uuid.Nil,
				Amount:   100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByUserID(gomock.Any(), inp.IDUser).Return(domain.User{}, sql.ErrNoRows)
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"could not find user: sql: no rows in result set"}`,
		}, {
			name:      "UserId_wallet_error",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","amount": 100}`,
			inputUser: wallet.InputDeposit{
				IDUser:   uuid.MustParse(IDUser),
				IDWallet: uuid.Nil,
				Amount:   100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByUserID(gomock.Any(), inp.IDUser).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{ID: uuid.Nil}}, nil)
				r.EXPECT().CreateWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet)}}, types.ErrorInternal)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"could not create wallet: internal error"}`,
		}, {
			name:      "email_user_not_found",
			inputBody: `{"email": "test@tgest.ru","amount": 100}`,
			inputUser: wallet.InputDeposit{
				EmailUser: "test@tgest.ru",
				IDUser:    uuid.Nil,
				IDWallet:  uuid.Nil,
				Amount:    100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByEmail(gomock.Any(), inp.EmailUser).Return(domain.User{}, sql.ErrNoRows)
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"could not find user: sql: no rows in result set"}`,
		}, {
			name:      "email_user_with_wallet",
			inputBody: `{"email": "test@tgest.ru","amount": 100}`,
			inputUser: wallet.InputDeposit{
				EmailUser: "test@tgest.ru",
				IDUser:    uuid.Nil,
				IDWallet:  uuid.Nil,
				Amount:    100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByEmail(gomock.Any(), inp.EmailUser).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.Nil}}, nil)
				r.EXPECT().CreateWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet)}}, types.ErrorInternal)

			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"could not create wallet: internal error"}`,
		}, {
			name:      "email_with_user_id",
			inputBody: `{"id_user": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","email": "test@tgest.ru","amount": 100}`,
			inputUser: wallet.InputDeposit{
				EmailUser: "test@tgest.ru",
				IDUser:    uuid.MustParse(IDUser),
				IDWallet:  uuid.Nil,
				Amount:    100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"too much user data provided, only one of the following is allowed: wallet, user_id, email"}`,
		}, {
			name:      "Only_amount",
			inputBody: `{"amount": 100}`,
			inputUser: wallet.InputDeposit{
				Amount: 100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"not enough user data provided"}`,
		}, {
			name:      "Internal_error_before_deposit",
			inputBody: `{"email": "test@tgest.ru","amount": 100}`,
			inputUser: wallet.InputDeposit{
				EmailUser: "test@tgest.ru",
				IDUser:    uuid.Nil,
				IDWallet:  uuid.Nil,
				Amount:    100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputDeposit) {
				r.EXPECT().CheckWalletByEmail(gomock.Any(), inp.EmailUser).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet)}}, nil)
				inp.IDWallet = uuid.MustParse(IDWallet)
				r.EXPECT().DepositWallet(gomock.Any(), inp).Return(domain.User{
					ID: uuid.MustParse(IDUser), Email: "test@tgest.ru", Wallet: domain.Wallet{
						ID: uuid.MustParse(IDWallet), Balance: 100}}, types.ErrorInternal)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"could not deposit to userinternal error"}`,
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
			walletRoutes.PUT("/deposit", h.DepositToUser)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/v1/user/wallet/deposit", bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			// Make Request
			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, strings.Split(w.Body.String(), "\n")[0], test.expectedResponseBody)
		})
	}
}

func TestHandler_TransferUsers(t *testing.T) {
	// Init Test Table

	type mockBehavior func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputTransferUsers)
	IDFrom := "7f6280a8-64fb-4408-b16b-1e77b9fa5f35"
	IDTo := "3ad280a8-64fb-4128-b16b-1e77b9f33f35"

	tests := []struct {
		name                 string
		inputBody            string
		input                wallet.InputTransferUsers
		balance              float32
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"from_id": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","to_id": "3ad280a8-64fb-4128-b16b-1e77b9f33f35","amount": 100}`,
			input: wallet.InputTransferUsers{
				FromID: uuid.MustParse(IDFrom),
				ToID:   uuid.MustParse(IDTo),
				Amount: 100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputTransferUsers) {
				r.EXPECT().CheckAndDoTransfer(gomock.Any(), inp).Return(domain.User{Email: "receiver@test.ru"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Transfered 100 to receiver@test.ru"}`,
		}, {
			name:      "One_id",
			inputBody: `{"from_id": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","to_id": "7f6280a8-64fb-4408-b16b-1e77b9fa5f35","amount": 100}`,
			input: wallet.InputTransferUsers{
				FromID: uuid.MustParse(IDFrom),
				ToID:   uuid.MustParse(IDFrom),
				Amount: 100,
			},
			balance: 100,
			mockBehavior: func(r *mocks_handler.MockUsers, ctx context.Context, inp wallet.InputTransferUsers) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"user cannot send money to himself"}`,
		},
	}

	// Run Tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mocks_handler.NewMockUsers(c)
			test.mockBehavior(users, context.Background(), test.input)
			h := Handler{usersService: users}
			router := echo.New()
			router.Validator = validator.NewValidator()

			v1 := router.Group("/v1")
			userRoutes := v1.Group("/user")
			userRoutes.POST("/", h.Create)
			// Init Endpoint
			walletRoutes := userRoutes.Group("/wallet")
			walletRoutes.PUT("/transfer", h.TransferUsers)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/v1/user/wallet/transfer", bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			// Make Request
			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, strings.Split(w.Body.String(), "\n")[0], test.expectedResponseBody)
		})
	}
}
