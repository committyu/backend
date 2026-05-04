package game

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
	"fmt"
)

type SyncGithubCommitUseCase struct {
	userRepo      repository.UserRepository
	gameDataRepo  repository.GameDataRepository
	githubService domain.GitHubService
}

func NewSyncGitHubCommitUsecase(
	userRepo repository.UserRepository,
	gameDataRepo repository.GameDataRepository,
	githubService domain.GitHubService,
) *SyncGithubCommitUseCase {
	return &SyncGithubCommitUseCase{
		userRepo:      userRepo,
		gameDataRepo:  gameDataRepo,
		githubService: githubService,
	}
}

func (u *SyncGithubCommitUseCase) Execute(
	ctx context.Context, userID domain.UserID,
) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	gameData, err := u.gameDataRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if gameData == nil {
		return fmt.Errorf("game data not found")
	}

	events, err := u.githubService.GetPushEvents(ctx, user.GithubName(), gameData.LastCommitCheckedAt())
	if err != nil {
		return err
	}

	newCommitCount := 0

	for _, event := range events {
		newCommitCount += event.CommitCount
	}

	if newCommitCount == 0 {
		return nil
	}

	gameData.AddCommits(newCommitCount)

	return u.gameDataRepo.Update(ctx, gameData)
}
