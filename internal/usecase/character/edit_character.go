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
	existingCharacter, err := c.characterRepo.FindByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingCharacter.UserID() != userID {
		return nil, domain.ErrCharacterNotFound
	}

	if update.Name != nil {
		existingCharacter.SetName(*update.Name)
	}
	if update.Hp != nil {
		existingCharacter.SetHp(*update.Hp)
	}
	if update.Atk != nil {
		existingCharacter.SetAtk(*update.Atk)
	}
	if update.Matk != nil {
		existingCharacter.SetMatk(*update.Matk)
	}
	if update.Def != nil {
		existingCharacter.SetDef(*update.Def)
	}
	if update.Mdef != nil {
		existingCharacter.SetMdef(*update.Mdef)
	}
	if update.Agi != nil {
		existingCharacter.SetAgi(*update.Agi)
	}
	if update.Luk != nil {
		existingCharacter.SetLuk(*update.Luk)
	}
	if update.Xp != nil {
		existingCharacter.SetXp(*update.Xp)
	}

	if err := c.characterRepo.Save(ctx, existingCharacter); err != nil {
		return nil, err
	}

	return existingCharacter, nil
}
