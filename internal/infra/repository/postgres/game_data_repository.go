package postgres

import (
	"context"
	"time"

	"backend/internal/domain"
	"backend/internal/infra/repository/model"
	"gorm.io/gorm"
)

type gameDataRepositoryImpl struct {
	db *gorm.DB
}

func NewGameDataRepository(db *gorm.DB) *gameDataRepositoryImpl {
	return &gameDataRepositoryImpl{db: db}
}

func (r *gameDataRepositoryImpl) Create(ctx context.Context, data *domain.GameData) error {
	m := model.GameData{
		UserID:              string(data.UserID()),
		MainCharacterID:     data.MainCharacterID(),
		PlayTime:            data.PlayTime(),
		Stage:               data.Stage(),
		LastCommitCheckedAt: data.LastCommitCheckedAt(),
		UpdatedAt:           time.Now(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *gameDataRepositoryImpl) Update(ctx context.Context, data *domain.GameData) error {
	m := model.GameData{
		UserID:              string(data.UserID()),
		MainCharacterID:     data.MainCharacterID(),
		PlayTime:            data.PlayTime(),
		Stage:               data.Stage(),
		LastCommitCheckedAt: data.LastCommitCheckedAt(),
		UpdatedAt:           time.Now(),
	}
	// UserIDをキーにして全フィールド更新
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *gameDataRepositoryImpl) FindByUserID(ctx context.Context, userID domain.UserID) (*domain.GameData, error) {
	var m model.GameData
	err := r.db.WithContext(ctx).Where("user_id = ?", string(userID)).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return domain.ReconstructGameData(
		domain.UserID(m.UserID),
		m.MainCharacterID,
		m.PlayTime,
		m.Stage,
		m.LastCommitCheckedAt,
		m.UpdatedAt,
	), nil
}