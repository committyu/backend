package postgres

import (
	"backend/internal/domain"
	"backend/internal/infra/repository/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db          *gorm.DB
	tokenCipher domain.TokenCipher
}

func NewUserRepository(db *gorm.DB, tokenCipher domain.TokenCipher) *userRepositoryImpl {
	return &userRepositoryImpl{
		db:          db,
		tokenCipher: tokenCipher,
	}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	encryptedToken, err := r.encryptToken(user.GithubAccessToken())
	if err != nil {
		return err
	}

	m := model.User{
		ID:                string(user.ID()),
		GithubName:        user.GithubName(),
		Email:             user.Email(),
		AvatarURL:         user.IconUrl(),
		GithubID:          user.GithubId(),
		GithubAccessToken: encryptedToken,
		CreatedAt:         user.CreatedAt(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *userRepositoryImpl) UpdateProfile(ctx context.Context, user *domain.User) error {
	encryptedToken, err := r.encryptToken(user.GithubAccessToken())
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", string(user.ID())).
		Updates(map[string]interface{}{
			"github_name":         user.GithubName(),
			"email":               user.Email(),
			"avatar_url":          user.IconUrl(),
			"github_access_token": encryptedToken,
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

	return r.toDomainUser(m)
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

	return r.toDomainUser(m)
}

func (r *userRepositoryImpl) encryptToken(token string) (string, error) {
	encryptedToken, err := r.tokenCipher.Encrypt(token)
	if err != nil {
		return "", fmt.Errorf("encrypt github access token failed: %w", err)
	}

	return encryptedToken, nil
}

func (r *userRepositoryImpl) decryptToken(token string) (string, error) {
	decryptedToken, err := r.tokenCipher.Decrypt(token)
	if err != nil {
		return "", fmt.Errorf("decrypt github access token failed: %w", err)
	}

	return decryptedToken, nil
}

func (r *userRepositoryImpl) toDomainUser(m model.User) (*domain.User, error) {
	githubAccessToken, err := r.decryptToken(m.GithubAccessToken)
	if err != nil {
		return nil, err
	}

	return domain.NewUser(
		domain.UserID(m.ID),
		m.GithubName,
		m.Email,
		m.AvatarURL,
		m.GithubID,
		githubAccessToken,
		m.CreatedAt,
	), nil
}
