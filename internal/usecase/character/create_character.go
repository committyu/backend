package character

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
)

type CreateCharacterUseCase struct {
	characterRepo repository.CharacterRepository
}

func NewCreateCharacterUseCase(cr repository.CharacterRepository) *CreateCharacterUseCase {
	return &CreateCharacterUseCase{
		characterRepo: cr,
	}
}

func (c *CreateCharacterUseCase) Execute(ctx context.Context, character *domain.Character) (*domain.Character, error) {
	createdCharacter, err := c.characterRepo.Create(ctx, character)
	if err != nil {
		return nil, err
	}

	return createdCharacter, nil
}
