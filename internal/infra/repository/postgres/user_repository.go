package postgres

import (
	"backend/internal/domain"
	"backend/internal/infra/repository/model"
	"context"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	m := model.User{
		ID:         string(user.ID()),
		GithubName: user.GithubName(),
		Email:      user.Email(),
		AvatarURL:  user.IconUrl(),
		GithubID:   user.GithubId(),
		CreatedAt:  user.CreatedAt(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *userRepositoryImpl) UpdateProfile(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", string(user.ID())).
		Updates(map[string]interface{}{
			"github_name": user.GithubName(),
			"email":       user.Email(),
			"avatar_url":  user.IconUrl(),
		}).Error
}

func (r *userRepositoryImpl) FindByGitHubID(ctx context.Context, githubID int64) (*domain.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).Where("github_id = ?", githubID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return domain.NewUser(
		domain.UserID(m.ID),
		m.GithubName,
		m.Email,
		m.AvatarURL,
		m.GithubID,
		m.CreatedAt,
	), nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var m model.User

	err := r.db.WithContext(ctx).Where("id = ?", string(id)).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return domain.NewUser(
		domain.UserID(m.ID),
		m.GithubName,
		m.Email,
		m.AvatarURL,
		m.GithubID,
		m.CreatedAt,
	), nil
}
