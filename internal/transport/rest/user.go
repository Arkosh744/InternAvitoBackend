package rest

import (
	"fmt"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Create UserAccount godoc
// @Summary     Create User
// @Description You need create account to use our service
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       userInputData body domain.InputUser true "Create User"
// @Router      /v1/user/ [post]
func (h *Handler) Create(ctx echo.Context) error {
	var userWallet domain.InputUser
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

// GetUserBalance Balance godoc
// @Summary     Get Data about User Balance
// @Description Data about User
// @Tags        account
// @Accept      json
// @Produce     json
// @Param       id path string true "User UUID"
// @Router      /v1/user/{id} [GET]
func (h *Handler) GetUserBalance(ctx echo.Context) error {
	var user domain.User
	var err error
	user.ID, err = uuid.Parse(ctx.Param("id"))
	log.Println(user.ID)
	log.Println(ctx.Request().RequestURI)
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
		"message": fmt.Sprintf("Balance of %s", user.Email),
		"balance": fmt.Sprintf("%v", user.Wallet.Balance),
	})
}
