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

func (r *characterRepositoryImpl) Edit(ctx context.Context, id domain.CharacterID, userID domain.UserID, update domain.CharacterUpdate) (*domain.Character, error) {
	updates := make(map[string]any)
	if update.Name != nil {
		updates["name"] = *update.Name
	}
	if update.Hp != nil {
		updates["hp"] = *update.Hp
	}
	if update.Atk != nil {
		updates["atk"] = *update.Atk
	}
	if update.Matk != nil {
		updates["matk"] = *update.Matk
	}
	if update.Def != nil {
		updates["def"] = *update.Def
	}
	if update.Mdef != nil {
		updates["mdef"] = *update.Mdef
	}
	if update.Agi != nil {
		updates["agi"] = *update.Agi
	}
	if update.Luk != nil {
		updates["luk"] = *update.Luk
	}
	if update.Xp != nil {
		updates["xp"] = *update.Xp
	}

	query := r.db.WithContext(ctx).Model(&model.Character{}).
		Where("id = ? AND user_id = ?", id.String(), userID.String())
	if len(updates) > 0 {
		result := query.Updates(updates)
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected == 0 {
			return nil, domain.ErrCharacterNotFound
		}
	}

	var m model.Character
	if err := query.First(&m).Error; err != nil {
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
