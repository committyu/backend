package repository

import (
	"backend/internal/domain"
	"context"
)

type CharacterRepository interface {
	Create(ctx context.Context, character *domain.Character) (*domain.Character, error)
	Edit(ctx context.Context, id domain.CharacterID, userID domain.UserID, update domain.CharacterUpdate) (*domain.Character, error)
}
