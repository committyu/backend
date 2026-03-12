package http

import (
	"errors"
	"net/http"

	"backend/internal/domain"
	"backend/internal/handler/http/presenter"
	"backend/internal/pkg/logger"
	"backend/internal/usecase/user"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUc *user.GetUserUsecase
}

func NewUserHandler(userUc *user.GetUserUsecase) *UserHandler {
	return &UserHandler{
		userUc: userUc,
	}
}
func (h *UserHandler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()

	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		logger.Error("invalid token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	u, err := h.userUc.Execute(ctx, domain.UserID(userIDStr))
	if err != nil {

		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "user not found",
			})
		}

		logger.Error("get_user usecase failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	res := presenter.ToUserResponse(u)

	logger.Info("user fetched success", "user_id", userIDStr)
	return c.JSON(http.StatusOK, res)
}