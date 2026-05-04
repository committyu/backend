package postgres

import (
	"context"
	"time"

	"backend/internal/domain"
	"backend/internal/infra/repository/model"
	"gorm.io/gorm"
)

type GameDataRepository struct {
	db *gorm.DB
}

func NewGameDataRepository(db *gorm.DB) *GameDataRepository {
	return &GameDataRepository{db: db}
}

func (r *GameDataRepository) Create(ctx context.Context, data *domain.GameData) error {
	m := model.GameData{
		UserID:              string(data.UserID()),
		MainCharacterID:     data.MainCharacterID(),
		PlayTime:            data.PlayTime(),
		Stage:               data.Stage(),
		GithubTotalCommits:  data.GithubTotalCommits(),
		LastCommitCheckedAt: data.LastCommitCheckedAt(),
		UpdatedAt:           time.Now(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *GameDataRepository) Update(ctx context.Context, data *domain.GameData) error {
	m := model.GameData{
		UserID:              string(data.UserID()),
		MainCharacterID:     data.MainCharacterID(),
		PlayTime:            data.PlayTime(),
		Stage:               data.Stage(),
		GithubTotalCommits:  data.GithubTotalCommits(),
		LastCommitCheckedAt: data.LastCommitCheckedAt(),
		UpdatedAt:           time.Now(),
	}
	// UserIDをキーにして全フィールド更新
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *GameDataRepository) FindByUserID(ctx context.Context, userID domain.UserID) (*domain.GameData, error) {
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
		m.GithubTotalCommits,
		m.LastCommitCheckedAt,
		m.UpdatedAt,
	), nil
}
