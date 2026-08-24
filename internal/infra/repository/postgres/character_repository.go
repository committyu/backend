package postgres

import (
	"backend/internal/domain"
	"backend/internal/infra/repository/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type characterRepositoryImpl struct {
	db *gorm.DB
}

func (r *characterRepositoryImpl) FindByCharacterID(ctx context.Context, id domain.CharacterID) (*domain.Character, error) {
	var m model.Character
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCharacterNotFound
		}
		return nil, err
	}
	characterID, err := domain.ParseCharacterID(m.ID)
	if err != nil {
		return nil, err
	}
	ownerID, err := domain.ParseUserID(m.UserID)
	if err != nil {
		return nil, err
	}
	return domain.RestoreCharacter(characterID, m.Name, m.Job, m.Hp, m.Atk, m.Matk, m.Def, m.Mdef, m.Agi, m.Luk, m.Xp, ownerID, m.CreatedAt), nil
}

func (r *characterRepositoryImpl) Save(ctx context.Context, character *domain.Character) error {
	result := r.db.WithContext(ctx).Model(&model.Character{}).
		Where("id = ? AND user_id = ?", character.ID().String(), character.UserID().String()).
		Updates(map[string]any{
			"name": character.Name(), "hp": character.Hp(), "atk": character.Atk(),
			"matk": character.Matk(), "def": character.Def(), "mdef": character.Mdef(),
			"agi": character.Agi(), "luk": character.Luk(), "xp": character.Xp(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrCharacterNotFound
	}
	return nil
}

func NewCharacterRepository(db *gorm.DB) *characterRepositoryImpl {
	return &characterRepositoryImpl{
		db: db,
	}
}

func (r *characterRepositoryImpl) Create(ctx context.Context, character *domain.Character) (*domain.Character, error) {
	m := model.Character{
		ID:        string(character.ID()),
		Name:      character.Name(),
		Job:       character.Job(),
		Hp:        character.Hp(),
		Atk:       character.Atk(),
		Matk:      character.Matk(),
		Def:       character.Def(),
		Mdef:      character.Mdef(),
		Agi:       character.Agi(),
		Luk:       character.Luk(),
		Xp:        character.Xp(),
		UserID:    string(character.UserID()),
		CreatedAt: character.CreatedAt(),
	}

	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}

	return character, nil
}
