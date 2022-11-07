package rest

import (
	"context"
	_ "github.com/Arkosh744/InternAvitoBackend/docs"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Users interface {
	Create(ctx context.Context, user domain.InputUser) (domain.User, error)
	GetUserBalance(ctx context.Context, user domain.User) (domain.User, error)
	CheckUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckWalletByUserID(ctx context.Context, uuid uuid.UUID) (domain.User, error)
	CheckWalletByEmail(ctx context.Context, user string) (domain.User, error)

	CreateWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error)
	DepositWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error)

	CheckAndDoTransfer(ctx context.Context, input wallet.InputTransferUsers) (domain.User, error)
	BuyServiceUser(ctx context.Context, input wallet.InputBuyServiceUser) (wallet.OutPendingOrder, error)
	ManageOrder(ctx context.Context, input wallet.InputOrderManager) (wallet.OutOrderManager, error)

	ReportMonth(ctx context.Context, input wallet.InputReportMonth) ([]wallet.ReportMonth, error)
	ReportForUser(ctx context.Context, input domain.InputReportUserTnx) ([]domain.OutputReportUserTnx, error)
}

type Handler struct {
	usersService Users
}

func NewHandler(users Users) *Handler {
	return &Handler{
		usersService: users,
	}
}

func (h *Handler) InitRouter() *echo.Echo {
	router := echo.New()
	router.Validator = validator.NewValidator()
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := router.Group("/v1")
	userRoutes := v1.Group("/user")
	userRoutes.POST("/", h.Create)
	userRoutes.GET("/:id", h.GetUserBalance)
	userRoutes.POST("/data", h.GetDataUser)

	walletRoutes := userRoutes.Group("/wallet")
	walletRoutes.PUT("/deposit", h.DepositToUser)
	// Когда-нибудь мы разрешим нашим пользователям выводить деньги, но не сегодня
	// walletRoutes.PUT("/withdrawal", h.Withdrawal)

	orderRoutes := walletRoutes.Group("/order")
	orderRoutes.PUT("/transfer", h.TransferUsers)
	orderRoutes.POST("/buy", h.BuyServiceUser)
	orderRoutes.POST("/approve", h.ManageOrder)
	orderRoutes.POST("/decline", h.ManageOrder)
	orderRoutes.POST("/report", h.ReportMonth)

	//router.Use(middleware.Logger())
	//router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//	Format: "method=${method}, uri=${uri}, status=${status}\n",
	//}))
	router.Use(middleware.Recover())

	return router
}
