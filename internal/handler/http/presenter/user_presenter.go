package presenter

import (
	"backend/internal/domain"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
}

func ToUserResponse(u *domain.User) UserResponse {
    return UserResponse{
        ID:        string(u.ID()),
		Name:      u.Name(),
		Email:     u.Email(),
		AvatarURL: u.IconUrl(),
		CreatedAt: u.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
    }
}