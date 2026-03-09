package router

import (
	"backend/internal/usecase/auth"
	"backend/internal/usecase/user"

	"github.com/labstack/echo/v4"
)

func StartEcho(
	loginUc *auth.LoginUsecase,
	tokenUc *auth.GenerateTokenUsecase,
	userUc *user.GetUserUsecase,
) {

	e := echo.New()

	api := e.Group("/api")

	RegisterAuthRoutes(api, loginUc, tokenUc)
	RegisterUserRoutes(api, userUc)

	e.Logger.Fatal(e.Start(":8080"))
}