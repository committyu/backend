package repository

import (
	"backend/internal/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByGitHubID(ctx context.Context, githubID int64) (*domain.User, error)
}