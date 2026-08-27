package character

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
)

type JobChangeCharacterUseCase struct {
	characterRepo repository.CharacterRepository
}

func NewJobChangeCharacterUseCase(cr repository.CharacterRepository) *JobChangeCharacterUseCase {
	return &JobChangeCharacterUseCase{
		characterRepo: cr,
	}
}

func (c *JobChangeCharacterUseCase) Execute(ctx context.Context, id domain.CharacterID, userID domain.UserID, job string) error {
	existingCharacter, err := c.characterRepo.FindByCharacterID(ctx, id)
	if err != nil {
		return err
	}
	if existingCharacter.UserID() != userID {
		return domain.ErrCharacterNotFound
	}

	err = c.characterRepo.JobChange(ctx, id, job)
	if err != nil {
		return err
	}

	return nil
}
