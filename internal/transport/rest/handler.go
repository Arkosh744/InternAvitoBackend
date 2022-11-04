package rest

import (
	"context"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/types"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Users interface {
	Create(ctx context.Context, user domain.WalletUser) (domain.WalletUser, error)
	CheckByEmail(ctx context.Context, email string) (domain.WalletUser, error)
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

	v1 := router.Group("/v1")
	userRoutes := v1.Group("/user")
	userRoutes.POST("/", h.Create)
	userRoutes.GET("/", h.List)
	userRoutes.GET("/:id", h.Get)
	userRoutes.PUT("/:id", h.Update)
	userRoutes.DELETE("/:id", h.Delete)

	// Middleware
	//router.Use(middleware.Logger())
	//router.Use(middleware.Recover())

	return router
}

func (h *Handler) Create(ctx echo.Context) error {
	var userWallet domain.WalletUser
	if err := ctx.Bind(&userWallet); err != nil {
		log.WithFields(log.Fields{"handler": "Create WalletUser"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := ctx.Validate(&userWallet); err != nil {
		log.WithFields(log.Fields{"handler": "Create WalletUser"}).Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}
	if _, err := h.usersService.CheckByEmail(ctx.Request().Context(), userWallet.Email); err != nil {
		log.WithFields(log.Fields{"handler": "Create WalletUser"}).Error(err)
		return ctx.JSON(http.StatusConflict, map[string]string{
			"message": "could not create user: " + err.Error(),
		})

	}

	createdUser, err := h.usersService.Create(ctx.Request().Context(), userWallet)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "could not create user" + err.Error(),
			})
		}
	}

	return ctx.JSON(http.StatusCreated, createdUser.ToWebUser())
}

func (h *Handler) List(ctx echo.Context) error {

	return nil
}

func (h *Handler) Get(ctx echo.Context) error {

	return nil
}

func (h *Handler) Update(ctx echo.Context) error {

	return nil
}

func (h *Handler) Delete(ctx echo.Context) error {

	return nil
}
