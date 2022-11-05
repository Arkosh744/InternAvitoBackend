package rest

import (
	"context"
	"fmt"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/types"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Users interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetUserBalance(ctx context.Context, user domain.User) (domain.User, error)
	CheckUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckWalletByUserID(ctx context.Context, uuid uuid.UUID) (domain.User, error)
	CheckWalletByEmail(ctx context.Context, user string) (domain.User, error)
	CreateWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error)
	DepositWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error)
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
	userRoutes.GET("/:id", h.GetUserBalance)

	walletRoutes := userRoutes.Group("/wallet")

	walletRoutes.PUT("/deposit", h.DepositToUser)
	//walletRoutes.PUT("/withdrawal", h.Update)
	//walletRoutes.PUT("/transfer", h.Update)
	//walletRoutes.DELETE("/:id", h.Delete)

	//router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	return router
}

func (h *Handler) Create(ctx echo.Context) error {
	var userWallet domain.User
	if err := ctx.Bind(&userWallet); err != nil {
		log.WithFields(log.Fields{"handler": "Create User"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := ctx.Validate(&userWallet); err != nil {
		log.WithFields(log.Fields{"handler": "Create User"}).Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}
	if _, err := h.usersService.CheckUserByEmail(ctx.Request().Context(), userWallet.Email); err != nil {
		log.WithFields(log.Fields{"handler": "Create User"}).Error(err)
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

func (h *Handler) DepositToUser(ctx echo.Context) error {
	var input wallet.InputDeposit
	if err := ctx.Bind(&input); err != nil {
		log.WithFields(log.Fields{"handler": "DepositToUser"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := ctx.Validate(&input); err != nil {
		log.WithFields(log.Fields{"handler": "DepositToUser"}).Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}

	if input.IDWallet != uuid.Nil && input.EmailUser == "" && input.IDUser == uuid.Nil {
		// just check it and continue
	} else if input.IDUser != uuid.Nil && input.EmailUser == "" && input.IDWallet == uuid.Nil {
		userWallet, err := h.usersService.CheckWalletByUserID(ctx.Request().Context(), input.IDUser)
		if err != nil {
			log.WithFields(log.Fields{"handler": "DepositToUser"}).Error(err)
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"message": "could not find user: " + err.Error(),
			})
		}
		if userWallet.Wallet.ID == uuid.Nil {
			userWallet, err = h.usersService.CreateWallet(ctx.Request().Context(), input)
			if err != nil {
				log.WithFields(log.Fields{"handler": "DepositToUser"}).Error(err)
				return ctx.JSON(http.StatusInternalServerError, map[string]string{
					"message": "could not create wallet: " + err.Error(),
				})
			}
		}
		input.IDWallet = userWallet.Wallet.ID
	} else if input.EmailUser != "" && input.IDUser == uuid.Nil && input.IDWallet == uuid.Nil {
		userWallet, err := h.usersService.CheckWalletByEmail(ctx.Request().Context(), input.EmailUser)
		if err != nil {
			log.WithFields(log.Fields{"handler": "DepositToUser"}).Error(err)
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"message": "could not find user: " + err.Error(),
			})
		}
		if userWallet.Wallet.ID == uuid.Nil {
			userWallet, err = h.usersService.CreateWallet(ctx.Request().Context(), input)
			if err != nil {
				log.WithFields(log.Fields{"handler": "DepositToUser"}).Error(err)
				return ctx.JSON(http.StatusInternalServerError, map[string]string{
					"message": "could not create wallet: " + err.Error(),
				})
			}
		}
		input.IDWallet = userWallet.Wallet.ID
	} else if input.IDWallet == uuid.Nil && input.EmailUser == "" && input.IDUser == uuid.Nil {
		log.WithFields(log.Fields{"handler": "DepositToUser"}).Error(types.ErrNotEnoughData)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": types.ErrNotEnoughData.Error(),
		})
	} else {
		log.WithFields(log.Fields{"handler": "DepositToUser"}).Error(types.ErrTooMuchData)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": types.ErrTooMuchData.Error(),
		})
	}
	user, err := h.usersService.DepositWallet(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "could not deposit to user" + err.Error(),
		})
	}
	log.Println(user)
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Deposited %v to %s. Balance: %v", input.Amount, user.Email, user.Wallet.Balance),
	})
}

func (h *Handler) GetUserBalance(ctx echo.Context) error {
	var user domain.User
	var err error
	user.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{"handler": "GetUserBalance"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	user, err = h.usersService.GetUserBalance(ctx.Request().Context(), user)
	if err == types.ErrNoWallet {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": err.Error(),
		})
	} else if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "could not get user balance " + err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("%s balance", user.Email),
		"balance": fmt.Sprintf("%v", user.Wallet.Balance),
	})
}
