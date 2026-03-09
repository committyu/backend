package router

import (
	"backend/internal/handler/http"
	"backend/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartEcho(authUc *usecase.AuthUsecase) {
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ハンドラーの初期化
	authHandler := http.NewAuthHandler(authUc)

	// ルーティング
	api := e.Group("/api/v1")
	{
		api.GET("/auth/callback", authHandler.Login)
	}

	// 起動
	e.Logger.Fatal(e.Start(":8080"))
}