package main

import (
	"os"

	infraAuth "backend/internal/infra/auth"
	infraCrypto "backend/internal/infra/crypto"
	"backend/internal/infra/db"
	"backend/internal/infra/github"
	"backend/internal/infra/repository/model"
	"backend/internal/infra/repository/postgres"
	"backend/internal/infra/router"
	"backend/internal/pkg/logger"
	"backend/internal/usecase/auth"
	"backend/internal/usecase/character"
	"backend/internal/usecase/game"
	"backend/internal/usecase/user"

	"github.com/joho/godotenv"
)

// @title Backend API
// @version 1.0
// @description GitHub OAuthとゲームデータを提供するAPIです。認証にはHttpOnly Cookieを使用します。
// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	logger.Init()
	if err := godotenv.Load(); err != nil {
		logger.Error(".env not found")
	}

	database, err := db.NewPostgresDB()
	if err != nil {
		logger.Error("database initialization failed", "error", err)
		os.Exit(1)
	}
	if err := database.AutoMigrate(&model.User{}, &model.GameData{}, &model.Character{}); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}

	tokenCipher, err := infraCrypto.NewTokenCipher(os.Getenv("GITHUB_TOKEN_ENCRYPTION_KEY"))
	if err != nil {
		logger.Error("token cipher initialization failed", "error", err)
		os.Exit(1)
	}
	userRepo := postgres.NewUserRepository(database, tokenCipher)
	gameRepo := postgres.NewGameDataRepository(database)
	characterRepo := postgres.NewCharacterRepository(database)

	githubClient := github.NewGitHubClient(github.Config{
		ClientID: os.Getenv("GITHUB_CLIENT_ID"), ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	})

	loginUc := auth.NewLoginUsecase(userRepo, gameRepo, githubClient)
	tokenUc := auth.NewGenerateTokenUsecase(infraAuth.NewJWTService())
	getUserUc := user.NewGetUserUsecase(userRepo)
	syncGithubCommitUc := game.NewSyncGitHubCommitUsecase(userRepo, gameRepo, githubClient)

	createCharacterUc := character.NewCreateCharacterUseCase(
		characterRepo,
	)
	editCharacterUc := character.NewEditCharacterUseCase(characterRepo)

	router.StartEcho(
		loginUc, tokenUc, getUserUc, syncGithubCommitUc, createCharacterUc, editCharacterUc,
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_REDIRECT_URL"),
		os.Getenv("FRONTEND_AUTH_CALLBACK_URL"),
	)
}
