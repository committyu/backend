package postgres

import (
	"context"
	"backend/internal/domain"
	"backend/internal/infra/repository/model" 
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
		ID:        string(user.ID()),
		Name:      user.Name(),
		Email:     user.Email(),
		AvatarURL: user.IconUrl(),
		GithubID:  user.GithubId(),
		CreatedAt: user.CreatedAt(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
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
		m.Name,
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
            return nil, nil // エラーファイル作ったら nil を変更　
        }
        return nil, err
    }

    return domain.NewUser(
        domain.UserID(m.ID),
        m.Name,
        m.Email,
        m.AvatarURL,
        m.GithubID,
        m.CreatedAt,
    ), nil
}