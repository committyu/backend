package usecase

import (
	"context"
	"fmt"
	"log"

	"backend/internal/domain"
	"backend/internal/repository"
)

// GitHubClientを抽象化したインターフェース
type GitHubService interface {
	GetUser(ctx context.Context, code string) (*domain.User, error)
}

type AuthUsecase struct {
	userRepo      repository.UserRepository
	gameDataRepo  repository.GameDataRepository
	githubService GitHubService
}

func NewAuthUsecase(
	ur repository.UserRepository,
	gr repository.GameDataRepository,
	gs GitHubService,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:      ur,
		gameDataRepo:  gr,
		githubService: gs,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, code string) (*domain.User, error) {
	// 1. GitHubからユーザー情報を取得（この時点のUserは一時的なもの）
	githubUser, err := u.githubService.GetUser(ctx, code)
	if err != nil {
		log.Printf("DEBUG: GetAccessToken failed: %v", err)
		return nil, fmt.Errorf("github auth failed: %w", err)
	}

	// 2. 既に登録されているかGitHub IDでチェック
	existingUser, err := u.userRepo.FindByGitHubID(ctx, githubUser.GithubId())
	if err == nil && existingUser != nil {
		// 既にいればそのユーザーを返す（ログイン成功）
		return existingUser, nil
	}

	// 3. いなければ新規登録処理
	// githubUserの情報を元に、永続化用のUserを再構築（IDなどは保持）
	// ここで domain.NewUser を呼び出し、一貫性を保つ
	newUser := githubUser 

	// ユーザーをDBに保存
	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 4. 初回ログインなので初期ゲームデータも作成
	newGameData := domain.NewGameData(newUser.ID())
	if err := u.gameDataRepo.Create(ctx, newGameData); err != nil {
		return nil, fmt.Errorf("failed to create game data: %w", err)
	}

	return newUser, nil
}