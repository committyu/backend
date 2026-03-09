package usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type GitHubService interface {
	GetUser(ctx context.Context, code string) (*domain.User, error)
}

type AuthUsecase struct {
	userRepo      repository.UserRepository
	gameDataRepo  repository.GameDataRepository
	githubService GitHubService
}

func NewAuthUsecase(ur repository.UserRepository, gr repository.GameDataRepository, gs GitHubService) *AuthUsecase {
	return &AuthUsecase{
		userRepo:     ur,
		gameDataRepo: gr,
		githubService: gs,
	}
}

// Login: GitHub認証とユーザー情報の永続化を行う
func (u *AuthUsecase) Login(ctx context.Context, code string) (*domain.User, error) {
	githubUser, err := u.githubService.GetUser(ctx, code)
	if err != nil {
		log.Printf("DEBUG: GitHub Auth failed: %v", err)
		return nil, fmt.Errorf("github auth failed: %w", err)
	}

	// 既存ユーザーチェック
	existingUser, err := u.userRepo.FindByGitHubID(ctx, githubUser.GithubId())
	if err == nil && existingUser != nil {
		return existingUser, nil
	}

	// 新規登録
	if err := u.userRepo.Create(ctx, githubUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 初期ゲームデータ作成
	newGameData := domain.NewGameData(githubUser.ID())
	if err := u.gameDataRepo.Create(ctx, newGameData); err != nil {
		return nil, fmt.Errorf("failed to create game data: %w", err)
	}

	return githubUser, nil
}

// GenerateToken: ユーザーIDを含むJWTを発行する
func (u *AuthUsecase) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID().String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 環境変数から取得。未設定なら安全のためエラーにする
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not defined in environment variables")
	}

	return token.SignedString([]byte(secret))
}