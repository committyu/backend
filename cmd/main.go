package main

import (
	"os"

	"backend/internal/infra/db"
	"backend/internal/infra/github"
	infraAuth "backend/internal/infra/auth"
	"backend/internal/infra/repository/model"
	"backend/internal/infra/repository/postgres"
	"backend/internal/infra/router"
	"backend/internal/usecase/auth"
	"backend/internal/usecase/user"
	"backend/internal/pkg/logger"

	"github.com/joho/godotenv"
)

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

	if err := database.AutoMigrate(
		&model.User{},
		&model.GameData{},
	); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}

	userRepo := postgres.NewUserRepository(database)
	gameRepo := postgres.NewGameDataRepository(database)

	githubClient := github.NewGitHubClient(github.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	})

	jwtService := infraAuth.NewJWTService()

	loginUc := auth.NewLoginUsecase(
		userRepo,
		gameRepo,
		githubClient,
	)

	tokenUc := auth.NewGenerateTokenUsecase(
		jwtService,
	)

	getUserUc := user.NewGetUserUsecase(
		userRepo,
	)

	router.StartEcho(
		loginUc,
		tokenUc,
		getUserUc,
	)
}