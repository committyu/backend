package http

import (
	"net/http"
	"time"

	"backend/internal/usecase/user"
	"backend/internal/domain"
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

	// ミドルウェアでセットしたuserIDを取得
	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	u, err := h.userUc.Execute(ctx, domain.UserID(userIDStr))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	type userResponse struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		AvatarURL string    `json:"avatar_url"`
		CreatedAt time.Time `json:"created_at"`
	}

	res := userResponse{
		ID:        string(u.ID()),
		Name:      u.Name(),
		Email:     u.Email(),
		AvatarURL: u.IconUrl(),
		CreatedAt: u.CreatedAt(),
	}

	return c.JSON(http.StatusOK, res)
}