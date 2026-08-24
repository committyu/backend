package character

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
)

type EditCharacterUseCase struct {
	characterRepo repository.CharacterRepository
}

func NewEditCharacterUseCase(cr repository.CharacterRepository) *EditCharacterUseCase {
	return &EditCharacterUseCase{
		characterRepo: cr,
	}
}

func (c *EditCharacterUseCase) Execute(ctx context.Context, id domain.CharacterID, userID domain.UserID, update domain.CharacterUpdate) (*domain.Character, error) {
	editedCharacter, err := c.characterRepo.Edit(ctx, id, userID, update)
	if err != nil {
		return nil, err
	}

	return editedCharacter, nil
}
