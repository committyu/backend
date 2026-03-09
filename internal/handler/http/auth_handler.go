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

func (h *AuthHandler) Login(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code is required"})
	}

	// 1. ログイン処理
	user, err := h.authUc.Login(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 2. JWT発行
	token, err := h.authUc.GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	// 3. 成功レスポンス
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}