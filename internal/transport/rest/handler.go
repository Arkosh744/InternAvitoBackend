package rest

import (
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/validator"
	"github.com/labstack/echo/v4"
)

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
	router.Validator = validator.NewValidator()

	v1 := router.Group("/v1")
	userRoutes := v1.Group("/user")
	userRoutes.POST("/", Create)
	userRoutes.GET("/:id", Get)
	userRoutes.GET("/", List)
	userRoutes.DELETE("/:id", Delete)
	//userRoutes.PUT("/:id", userController.Update)

	return router
}
