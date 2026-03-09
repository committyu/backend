package repository 

import (
    "context"
	"backend/internal/domain"
)

type GameDataRepository interface {
    Create(ctx context.Context, data *domain.GameData) error
    Update(ctx context.Context, data *domain.GameData) error
}