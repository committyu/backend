package main

import (
	"backend/internal/infra/db"
	"backend/internal/infra/github"
	"backend/internal/infra/repository/postgres"
	"backend/internal/infra/router"
	"backend/internal/usecase"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load(".env")

    db, err := db.NewPostgresDB()
    if err != nil {
        log.Fatal(err)
    }

    userRepo := postgres.NewUserRepository(db)
    gameRepo := postgres.NewGameDataRepository(db)
    ghClient := github.NewGitHubClient(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"))

    authUc := usecase.NewAuthUsecase(userRepo, gameRepo, ghClient)

    router.StartEcho(authUc)
}