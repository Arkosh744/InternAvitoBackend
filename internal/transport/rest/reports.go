package rest

import (
	"fmt"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// ReportMonth Report godoc
// @Summary     Preparing a monthly report for the accounting department
// @Description Preparing a monthly report for the accounting department
// @Tags        Reports
// @Accept      json
// @Produce     json
// @Param       InputDate body wallet.InputReportMonth true "Input month and year"
// @Router      /v1/user/wallet/order/report [POST]
func (h *Handler) ReportMonth(ctx echo.Context) error {
	var input wallet.InputReportMonth
	if err := ctx.Bind(&input); err != nil {
		log.WithFields(log.Fields{"handler": "ReportDataUser"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := ctx.Validate(&input); err != nil {
		log.WithFields(log.Fields{"handler": "ReportDataUser"}).Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}

	report, err := h.usersService.ReportMonth(ctx.Request().Context(), input)
	if err != nil {
		log.WithFields(log.Fields{"handler": "ReportDataUser"}).Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"message": fmt.Sprintf("Report for %v %v", input.Month, input.Year),
		"report":  report,
	})
}

// GetDataUser Report godoc
// @Summary     Get data about transactions for User
// @Description Pagination is available
// @Tags        Reports
// @Accept      json
// @Produce     json
// @Param       UserDataAndLimits body domain.InputReportUserTnx true "Input User Data And Limits"
// @Router      /v1/user/data [POST]
func (h *Handler) GetDataUser(ctx echo.Context) error {
	var input domain.InputReportUserTnx
	if err := ctx.Bind(&input); err != nil {
		log.WithFields(log.Fields{"handler": "ReportDataUser"}).Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := ctx.Validate(&input); err != nil {
		log.WithFields(log.Fields{"handler": "ReportDataUser"}).Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}

	report, err := h.usersService.ReportForUser(ctx.Request().Context(), input)
	if err != nil {
		log.WithFields(log.Fields{"handler": "ReportDataUser"}).Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"message": fmt.Sprintf("Report for %s", input.IDUser),
		"report":  report,
	})

}
