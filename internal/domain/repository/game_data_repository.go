package repository

import (
	"backend/internal/domain"
	"context"
)

type GameDataRepository interface {
	Create(ctx context.Context, data *domain.GameData) error
	Update(ctx context.Context, data *domain.GameData) error
	FindByUserID(ctx context.Context, userID domain.UserID) (*domain.GameData, error)
	ResetPlayTime(ctx context.Context) error
}