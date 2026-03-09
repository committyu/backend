package main

import (
	"log"
	"os"

	"backend/internal/infra/db"
	"backend/internal/infra/github"
	infraAuth "backend/internal/infra/auth"
	"backend/internal/infra/repository/model"
	"backend/internal/infra/repository/postgres"
	"backend/internal/infra/router"
	"backend/internal/usecase/auth"
	"backend/internal/usecase/user"

	"github.com/joho/godotenv"
)

func main() {

	// env
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	// DB
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	// migration
	if err := database.AutoMigrate(
		&model.User{},
		&model.GameData{},
	); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// repositories
	userRepo := postgres.NewUserRepository(database)
	gameRepo := postgres.NewGameDataRepository(database)

	// external services
	githubClient := github.NewGitHubClient(github.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	})

	jwtService := infraAuth.NewJWTService()

	// usecases
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

	// start server
	router.StartEcho(
		loginUc,
		tokenUc,
		getUserUc,
	)
}