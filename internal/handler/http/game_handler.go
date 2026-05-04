package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"backend/internal/domain"
	"backend/internal/usecase/game"
)

type GameHandler struct {
	syncUseCase *game.SyncGithubCommitUseCase
}

func NewGameHandler(syncUseCase *game.SyncGithubCommitUseCase) *GameHandler {
	return &GameHandler{
		syncUseCase: syncUseCase,
	}
}

func (h *GameHandler) SyncGithubCommits(c echo.Context) error {
	userIDStr, ok := c.Get("userID").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token",
		})
	}

	userID, err := domain.ParseUserID(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid user id",
		})
	}

	err = h.syncUseCase.Execute(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "sync completed",
	})
}
