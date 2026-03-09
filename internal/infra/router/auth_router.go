package router

import (
	"backend/internal/handler/http"
	"backend/internal/usecase/auth"

	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(
	api *echo.Group,
	loginUc *auth.LoginUsecase,
	tokenUc *auth.GenerateTokenUsecase,
) {

	authHandler := http.NewAuthHandler(loginUc, tokenUc)

	auth := api.Group("/auth")
	{
		auth.GET("/callback", authHandler.Login)
	}
}