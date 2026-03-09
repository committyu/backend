package domain

import "context"

type GitHubService interface {
	GetUser(ctx context.Context, code string) (*User, error)
}