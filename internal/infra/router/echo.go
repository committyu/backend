package router

import (
	_ "backend/docs"
	"backend/internal/pkg/logger"
	"backend/internal/usecase/auth"
	"backend/internal/usecase/game"
	"backend/internal/usecase/user"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func StartEcho(
	loginUc *auth.LoginUsecase,
	tokenUc *auth.GenerateTokenUsecase,
	userUc *user.GetUserUsecase,
	syncGithubCommitUc *game.SyncGithubCommitUseCase,
) {

	e := echo.New()

	e.Use(logger.RequestLogger())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api")

	RegisterAuthRoutes(api, loginUc, tokenUc)
	RegisterUserRoutes(api, userUc)
	RegisterGameRoutes(api, syncGithubCommitUc)

	e.Logger.Fatal(e.Start(":8080"))
}
