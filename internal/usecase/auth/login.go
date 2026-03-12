package auth

import (
	"context"
	"fmt"

	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/pkg/logger"
)

type LoginUsecase struct {
	userRepo      repository.UserRepository
	gameDataRepo  repository.GameDataRepository
	githubService domain.GitHubService
}

func NewLoginUsecase(
	ur repository.UserRepository,
	gr repository.GameDataRepository,
	gs domain.GitHubService,
) *LoginUsecase {

	return &LoginUsecase{
		userRepo:      ur,
		gameDataRepo:  gr,
		githubService: gs,
	}
}

func (u *LoginUsecase) Execute(
	ctx context.Context,
	code string,
) (*domain.User, error) {

	githubUser, err := u.githubService.GetUser(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("github auth failed: %w", err)
	}

	existingUser, err := u.userRepo.FindByGitHubID(
		ctx,
		githubUser.GithubId(),
	)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		logger.Info("user login", "user_id", existingUser.ID())
		return existingUser, nil
	}

	logger.Info("new user signup", "github_id", githubUser.GithubId())

	if err := u.userRepo.Create(ctx, githubUser); err != nil {
		return nil, err
	}

	gameData := domain.NewGameData(githubUser.ID())

	if err := u.gameDataRepo.Create(ctx, gameData); err != nil {
		return nil, err
	}

	return githubUser, nil
}