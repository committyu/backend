package http

import (
	"backend/internal/domain"
	"backend/internal/handler/http/presenter"
	"backend/internal/usecase/character"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type CharacterHandler struct {
	createUc *character.CreateCharacterUseCase
}

func NewCharacterHandler(createUc *character.CreateCharacterUseCase) *CharacterHandler {
	return &CharacterHandler{
		createUc: createUc,
	}
}

func (h *CharacterHandler) Create(c echo.Context) error {
	var req presenter.CreateCharacterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	userIDValue, ok := c.Get("userID").(string)
	if !ok || userIDValue == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	userID, err := domain.ParseUserID(userIDValue)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	newCharacter := domain.NewCharacter(
		req.Name,
		req.Job,
		userID,
		time.Now(),
	)

	createdCharacter, err := h.createUc.Execute(c.Request().Context(), newCharacter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create character"})
	}

	return c.JSON(http.StatusCreated, presenter.CreateCharacterResponse{
		ID:        createdCharacter.ID().String(),
		Name:      createdCharacter.Name(),
		Job:       createdCharacter.Job(),
		UserID:    createdCharacter.UserID().String(),
		CreatedAt: createdCharacter.CreatedAt(),
	})
}
