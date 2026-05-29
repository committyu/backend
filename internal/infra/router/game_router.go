package router

import (
	"backend/internal/handler/http"
	"backend/internal/infra/middleware"
	"backend/internal/usecase/game"

	"github.com/labstack/echo/v4"
)

func RegisterGameRoutes(
	api *echo.Group,
	syncGithubCommitUc *game.SyncGithubCommitUseCase,
) {
	gameHandler := http.NewGameHandler(syncGithubCommitUc)

	game := api.Group("/game")
	game.Use(middleware.AuthMiddleware())

	{
		game.GET("/syncCommit", gameHandler.SyncGithubCommits)
	}
}
