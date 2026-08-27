package character

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
)

type StatusEditCharacterUseCase struct {
	characterRepo repository.CharacterRepository
}

type StatusUpdate struct {
	Hp   *int
	Atk  *int
	Matk *int
	Def  *int
	Mdef *int
	Agi  *int
	Luk  *int
}

func NewStatusEditCharacterUseCase(cr repository.CharacterRepository) *StatusEditCharacterUseCase {
	return &StatusEditCharacterUseCase{
		characterRepo: cr,
	}
}

func (c *StatusEditCharacterUseCase) Execute(ctx context.Context, id domain.CharacterID, userID domain.UserID, xp int, update StatusUpdate) (*domain.Character, error) {
	if xp <= 0 {
		return nil, domain.ErrInvalidXP
	}
	if !validStatusUpdate(xp, update) {
		return nil, domain.ErrInvalidStatus
	}

	existingCharacter, err := c.characterRepo.FindByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingCharacter.UserID() != userID {
		return nil, domain.ErrCharacterNotFound
	}

	if existingCharacter.Xp() < xp {
		return nil, domain.ErrInsufficientXP
	}

	if update.Hp != nil {
		existingCharacter.SetHp(existingCharacter.Hp() + *update.Hp)
	}
	if update.Atk != nil {
		existingCharacter.SetAtk(existingCharacter.Atk() + *update.Atk)
	}
	if update.Matk != nil {
		existingCharacter.SetMatk(existingCharacter.Matk() + *update.Matk)
	}
	if update.Def != nil {
		existingCharacter.SetDef(existingCharacter.Def() + *update.Def)
	}
	if update.Mdef != nil {
		existingCharacter.SetMdef(existingCharacter.Mdef() + *update.Mdef)
	}
	if update.Agi != nil {
		existingCharacter.SetAgi(existingCharacter.Agi() + *update.Agi)
	}
	if update.Luk != nil {
		existingCharacter.SetLuk(existingCharacter.Luk() + *update.Luk)
	}
	existingCharacter.SetXp(existingCharacter.Xp() - xp)

	if err := c.characterRepo.Save(ctx, existingCharacter); err != nil {
		return nil, err
	}

	return existingCharacter, nil
}

func validStatusUpdate(xp int, update StatusUpdate) bool {
	values := []*int{update.Hp, update.Atk, update.Matk, update.Def, update.Mdef, update.Agi, update.Luk}
	count := 0
	total := 0
	for _, value := range values {
		if value != nil {
			if *value <= 0 {
				return false
			}
			count++
			total += *value
		}
	}
	return count > 0 && total == xp
}
