package rest

import (
	"fmt"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// DepositToUser Deposit godoc
// @Summary     Deposit funds to user
// @Description We receive funds, now process it to user wallet. If wallet doesn't exist, we create it
// @Description You must provide only one of the fields: user_id or email or wallet_id AND amount
// @Tags        wallet
// @Accept      json
// @Produce     json
// @Param       DepositData body wallet.InputDeposit true "You must provide only one of the user fields and amount: (user_id || email || wallet_id) && amount"
// @Router      /v1/user/wallet/deposit [put]
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
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Deposited %v to %s. Balance: %v", input.Amount, user.Email, user.Wallet.Balance),
	})
}

// TransferUsers Transfer godoc
// @Summary     Transfer funds between 2 Users
// @Description Transfer funds from 1 user to another
// @Tags        wallet
// @Accept      json
// @Produce     json
// @Param       TrasferInput body wallet.InputTransferUsers true "Trasfer Input between 2 Users IDs"
// @Router      /v1/user/wallet/ [PUT]
func (h *Handler) TransferUsers(ctx echo.Context) error {
	var input wallet.InputTransferUsers
	if err := ctx.Bind(&input); err != nil {
		log.WithFields(log.Fields{"handler": "TransferUsers"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := ctx.Validate(&input); err != nil {
		log.WithFields(log.Fields{"handler": "TransferUsers"}).Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}
	if input.FromID == input.ToID {
		log.WithFields(log.Fields{"handler": "TransferUsers"}).Error(types.ErrSameUser)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": types.ErrSameUser.Error(),
		})
	}

	toUser, err := h.usersService.CheckAndDoTransfer(ctx.Request().Context(), input)
	if err != nil {
		log.WithFields(log.Fields{"handler": "TransferUsers"}).Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Transfered %v to %s", input.Amount, toUser.Email),
	})
}
