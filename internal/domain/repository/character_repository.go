package repository

import (
	"backend/internal/domain"
	"context"
)

type CharacterRepository interface {
	Create(ctx context.Context, character *domain.Character) (*domain.Character, error)
	FindByCharacterID(ctx context.Context, id domain.CharacterID) (*domain.Character, error)
	Save(ctx context.Context, character *domain.Character) error
}
