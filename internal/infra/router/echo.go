package router

import (
	"backend/internal/pkg/logger"
	"backend/internal/usecase/auth"
	"backend/internal/usecase/game"
	"backend/internal/usecase/user"

	"github.com/labstack/echo/v4"
)

func StartEcho(
	loginUc *auth.LoginUsecase,
	tokenUc *auth.GenerateTokenUsecase,
	userUc *user.GetUserUsecase,
	syncGithubCommitUc *game.SyncGithubCommitUseCase,
) {

	e := echo.New()

	e.Use(logger.RequestLogger())

	api := e.Group("/api")

	RegisterAuthRoutes(api, loginUc, tokenUc)
	RegisterUserRoutes(api, userUc)
	RegisterGameRoutes(api, syncGithubCommitUc)

	e.Logger.Fatal(e.Start(":8080"))
}
