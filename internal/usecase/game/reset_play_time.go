package game

import (
	"backend/internal/domain/repository"
	"context"
)

type ResetPlayTime struct {
	repo repository.GameDataRepository
}

func NewResetPlayTime(repo repository.GameDataRepository) *ResetPlayTime {
	return &ResetPlayTime{repo: repo}
}

func (u *ResetPlayTime) Execute(ctx context.Context) error {
	return u.repo.ResetPlayTime(ctx)
}