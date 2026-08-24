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
	githubClientID string,
	githubRedirectURL string,
	frontendCallbackURL string,
) {

	authHandler := http.NewAuthHandler(loginUc, tokenUc, githubClientID, githubRedirectURL, frontendCallbackURL)

	auth := api.Group("/auth")
	{
		auth.GET("/login", authHandler.Login)
		auth.GET("/callback", authHandler.Callback)
	}
}
