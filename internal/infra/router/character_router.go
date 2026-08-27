package router

import (
	"backend/internal/handler/http"
	"backend/internal/infra/middleware"
	"backend/internal/usecase/character"

	"github.com/labstack/echo/v4"
)

func RegisterCharacterRoutes(
	api *echo.Group,
	createUc *character.CreateCharacterUseCase,
	statusEditUc *character.StatusEditCharacterUseCase,
	jobChangeUc *character.JobChangeCharacterUseCase,
) {

	characterHandler := http.NewCharacterHandler(createUc, statusEditUc, jobChangeUc)

	character := api.Group("/character")
	character.Use(middleware.AuthMiddleware())
	{
		character.POST("", characterHandler.Create)
		character.PATCH("/status", characterHandler.StatusEdit)
		character.PATCH("/job", characterHandler.JobChange)
	}
}
