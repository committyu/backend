package game

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"context"
	"fmt"
	"time"
)

type SyncGithubCommitResult struct {
	GithubName            string
	CheckedSince          time.Time
	CheckedUntil          time.Time
	MatchedPushEventCount int
	NewCommitCount        int
	TotalCommits          int
	Updated               bool
}

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
) (*SyncGithubCommitResult, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	gameData, err := u.gameDataRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if gameData == nil {
		return nil, fmt.Errorf("game data not found")
	}

	checkedSince := gameData.LastCommitCheckedAt()
	events, err := u.githubService.GetPushEvents(ctx, user.GithubName(), checkedSince)
	if err != nil {
		return nil, err
	}

	newCommitCount := 0

	for _, event := range events {
		newCommitCount += event.CommitCount
	}

	result := &SyncGithubCommitResult{
		GithubName:            user.GithubName(),
		CheckedSince:          checkedSince,
		CheckedUntil:          time.Now(),
		MatchedPushEventCount: len(events),
		NewCommitCount:        newCommitCount,
		TotalCommits:          gameData.GithubTotalCommits(),
		Updated:               false,
	}

	if newCommitCount == 0 {
		return result, nil
	}

	gameData.AddCommits(newCommitCount)
	result.TotalCommits = gameData.GithubTotalCommits()
	result.Updated = true

	return result, u.gameDataRepo.Update(ctx, gameData)
}
