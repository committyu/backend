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

// Create: 初回保存
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

// Update: プレイ時間やステージの更新
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

// FindByUserID: ユーザーIDで検索してドメインモデルに変換
func (r *gameDataRepositoryImpl) FindByUserID(ctx context.Context, userID domain.UserID) (*domain.GameData, error) {
	var m model.GameData
	err := r.db.WithContext(ctx).Where("user_id = ?", string(userID)).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// ドメインモデルの再構築
	// ※domain.NewGameDataは初期値用なので、取得した値をセットするロジックが必要です
	data := domain.NewGameData(domain.UserID(m.UserID))
	
	// 注意: domain.GameDataのフィールドが非公開なので、
	// domain側に「DBからの復元用関数」か「セッター」を用意して反映させるのが一般的です
	
	return data, nil
}