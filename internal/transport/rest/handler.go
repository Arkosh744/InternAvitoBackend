package rest

import "github.com/labstack/echo/v4"

type Users interface {
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
	return router
}
