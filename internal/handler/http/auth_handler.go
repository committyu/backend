package http

import (
	"net/http"

	"backend/internal/pkg/logger"
	"backend/internal/usecase/auth"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	loginUc *auth.LoginUsecase
	tokenUc *auth.GenerateTokenUsecase
}

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

func NewAuthHandler(
	loginUc *auth.LoginUsecase,
	tokenUc *auth.GenerateTokenUsecase,
) *AuthHandler {
	return &AuthHandler{
		loginUc: loginUc,
		tokenUc: tokenUc,
	}
}

// Login godoc
// @Summary GitHub OAuth でログイン
// @Description GitHub OAuth の一時コードを JWT に交換します。
// @Tags auth
// @Produce json
// @Param code query string true "GitHub OAuth code"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [get]
// @Router /auth/callback [get]
func (h *AuthHandler) Login(c echo.Context) error {

	ctx := c.Request().Context()

	code := c.QueryParam("code")
	if code == "" {
		logger.Error("github oauth code missing")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "code is required",
		})
	}

	user, err := h.loginUc.Execute(ctx, code)
	if err != nil {
		logger.Error("login usecase failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	token, err := h.tokenUc.Execute(user.ID())
	if err != nil {
		logger.Error("token generation failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate token",
		})
	}

	logger.Info("login success", "user_id", user.ID())
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
