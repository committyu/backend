package domain

import (
	"context"
	"time"
)

type GitHubPushEvent struct {
	ID          string
	CreatedAt   time.Time
	CommitCount int
}

type GitHubService interface {
	GetUser(ctx context.Context, code string) (*User, error)

	GetPushEvents(ctx context.Context, username string, accessToken string, lastCommitCheckedAt time.Time) ([]GitHubPushEvent, error)
}
