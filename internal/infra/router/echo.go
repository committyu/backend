package router

import (
	"net/http"
	"net/url"

	_ "backend/docs"
	"backend/internal/pkg/logger"
	"backend/internal/usecase/auth"
	"backend/internal/usecase/character"
	"backend/internal/usecase/game"
	"backend/internal/usecase/user"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func StartEcho(
	loginUc *auth.LoginUsecase,
	tokenUc *auth.GenerateTokenUsecase,
	userUc *user.GetUserUsecase,
	syncGithubCommitUc *game.SyncGithubCommitUseCase,
	characterCreateUc *character.CreateCharacterUseCase,
	characterEditUc *character.EditCharacterUseCase,
	githubClientID string,
	githubRedirectURL string,
	frontendCallbackURL string,
) {

	e := echo.New()

	e.Use(logger.RequestLogger())
	if frontendURL, err := url.Parse(frontendCallbackURL); err == nil && frontendURL.Scheme != "" && frontendURL.Host != "" {
		e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins:     []string{frontendURL.Scheme + "://" + frontendURL.Host},
			AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.OPTIONS},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowCredentials: true,
		}))
	}
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(
			http.StatusTemporaryRedirect,
			"/swagger/index.html",
		)
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api")

	RegisterAuthRoutes(api, loginUc, tokenUc, githubClientID, githubRedirectURL, frontendCallbackURL)
	RegisterUserRoutes(api, userUc)
	RegisterGameRoutes(api, syncGithubCommitUc)
	RegisterCharacterRoutes(api, characterCreateUc, characterEditUc)

	e.Logger.Fatal(e.Start(":8080"))
}
