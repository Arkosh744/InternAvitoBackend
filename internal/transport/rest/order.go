package rest

import (
	"fmt"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// BuyServiceUser Services godoc
// @Summary     Buy Service for User from another -> make "pending" order
// @Description Buy Service for User from another -> make "pending" order that should be proceeding with /v1/user/wallet/approve or /v1/user/wallet/decline
// @Description Initially we have only 3 services: Dodo Pizza, Yandex Taxi, Yandex Food
// @Tags        order
// @Accept      json
// @Produce     json
// @Param       InputBuy body wallet.InputBuyServiceUser true "Buy Service Input"
// @Router      /v1/user/wallet/buy [POST]
func (h *Handler) BuyServiceUser(ctx echo.Context) error {
	var input wallet.InputBuyServiceUser
	if err := ctx.Bind(&input); err != nil {
		log.WithFields(log.Fields{"handler": "BuyServiceUser"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := ctx.Validate(&input); err != nil {
		log.WithFields(log.Fields{"handler": "BuyServiceUser"}).Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}

	order, err := h.usersService.BuyServiceUser(ctx.Request().Context(), input)
	if err != nil {
		log.WithFields(log.Fields{"handler": "BuyServiceUser"}).Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"message":  fmt.Sprintf("Created order for %s for %v", order.ServiceName, order.Cost),
		"order_id": order.ID.String(),
	})
}

// ManageOrder Manage godoc
// @Summary     Approve or Decline Order for User depends on endpoint
// @Description Approve or Decline order
// @Tags        order
// @Accept      json
// @Produce     json
// @Param       InputBuy body wallet.InputBuyServiceUser true "Buy Service Input"
// @Router      /v1/user/wallet/order/approve [POST]
// @Router      /v1/user/wallet/order/decline [POST]
func (h *Handler) ManageOrder(ctx echo.Context) error {
	var input wallet.InputOrderManager
	if err := ctx.Bind(&input); err != nil {
		log.WithFields(log.Fields{"handler": "ApproveOrder"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := ctx.Validate(&input); err != nil {
		log.WithFields(log.Fields{"handler": "ApproveOrder"}).Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}
	reqUri := strings.Split(ctx.Request().RequestURI, "order/")[1]

	if reqUri == "approve" {
		input.Status = "approved"
	} else if reqUri == "decline" {
		input.Status = "cancelled"
	}
	outOrder, err := h.usersService.ManageOrder(ctx.Request().Context(), input)
	if err != nil {
		log.WithFields(log.Fields{"handler": "ApproveOrder"}).Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Order %s", outOrder.Status),
	})
}
