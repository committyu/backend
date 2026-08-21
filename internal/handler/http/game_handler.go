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

type SyncGithubCommitsResponse struct {
	Message               string `json:"message"`
	GithubName            string `json:"github_name"`
	CheckedSince          string `json:"checked_since"`
	CheckedUntil          string `json:"checked_until"`
	MatchedPushEventCount int    `json:"matched_push_event_count"`
	NewCommitCount        int    `json:"new_commit_count"`
	TotalCommits          int    `json:"total_commits"`
	Updated               bool   `json:"updated"`
}

func NewGameHandler(syncUseCase *game.SyncGithubCommitUseCase) *GameHandler {
	return &GameHandler{
		syncUseCase: syncUseCase,
	}
}

// SyncGithubCommits godoc
// @Summary GitHub のコミット数を同期
// @Description GitHub の PushEvent を取得し、ゲームデータのコミット数を更新します。
// @Tags game
// @Produce json
// @Success 200 {object} SyncGithubCommitsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /game/syncCommit [get]
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

	result, err := h.syncUseCase.Execute(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SyncGithubCommitsResponse{
		Message:               "sync completed",
		GithubName:            result.GithubName,
		CheckedSince:          result.CheckedSince.Format("2006-01-02T15:04:05Z07:00"),
		CheckedUntil:          result.CheckedUntil.Format("2006-01-02T15:04:05Z07:00"),
		MatchedPushEventCount: result.MatchedPushEventCount,
		NewCommitCount:        result.NewCommitCount,
		TotalCommits:          result.TotalCommits,
		Updated:               result.Updated,
	})
}
