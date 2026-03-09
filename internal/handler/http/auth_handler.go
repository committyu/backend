package http

import (
	"net/http"

	"backend/internal/usecase/auth"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	loginUc *auth.LoginUsecase
	tokenUc *auth.GenerateTokenUsecase
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

func (h *AuthHandler) Login(c echo.Context) error {

	ctx := c.Request().Context()

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "code is required",
		})
	}

	user, err := h.loginUc.Execute(ctx, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	token, err := h.tokenUc.Execute(user.ID())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}