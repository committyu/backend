package http

import (
	"net/http"
	"backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUc *usecase.AuthUsecase
}

func NewAuthHandler(authUc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUc: authUc}
}

// Login: GitHubのコールバックを受けてログイン処理を行う
func (h *AuthHandler) Login(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code is required"})
	}

	user, err := h.authUc.Login(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}