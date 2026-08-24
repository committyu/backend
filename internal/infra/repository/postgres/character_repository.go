package postgres

import (
	"backend/internal/domain"
	"backend/internal/infra/repository/model"
	"context"

	"gorm.io/gorm"
)

type characterRepositoryImpl struct {
	db *gorm.DB
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
