package repository

import (
	"backend/internal/domain"
	"context"
)

type CharacterRepository interface {
	Create(ctx context.Context, character *domain.Character) (*domain.Character, error)
}