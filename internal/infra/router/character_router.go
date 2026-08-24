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
	editUc *character.EditCharacterUseCase,
) {

	characterHandler := http.NewCharacterHandler(createUc, editUc)

	character := api.Group("/character")
	character.Use(middleware.AuthMiddleware())
	{
		character.POST("", characterHandler.Create)
		character.PUT("", characterHandler.Edit)
	}
}
