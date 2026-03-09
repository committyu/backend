package router

import (
	"backend/internal/handler/http"
	"backend/internal/infra/middleware"
	"backend/internal/usecase/user"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(
	api *echo.Group,
	getUserUc *user.GetUserUsecase,
) {

	userHandler := http.NewUserHandler(getUserUc)

	user := api.Group("/user")

    user.Use(middleware.AuthMiddleware())

	{
		user.GET("/me", userHandler.GetMe)
	}
}