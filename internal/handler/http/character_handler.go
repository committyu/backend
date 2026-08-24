package http

import (
	"backend/internal/domain"
	"backend/internal/handler/http/presenter"
	"backend/internal/usecase/character"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type CharacterHandler struct {
	createUc *character.CreateCharacterUseCase
	editUc   *character.EditCharacterUseCase
}

func NewCharacterHandler(createUc *character.CreateCharacterUseCase, editUc *character.EditCharacterUseCase) *CharacterHandler {
	return &CharacterHandler{
		createUc: createUc,
		editUc:   editUc,
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

func (h *CharacterHandler) Edit(c echo.Context) error {
	var req presenter.EditCharacterReq
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

	characterID, err := domain.ParseCharacterID(req.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid character id"})
	}

	editedCharacter, err := h.editUc.Execute(c.Request().Context(), characterID, userID, domain.CharacterUpdate{
		Name: req.Name, Hp: req.Hp, Atk: req.Atk, Matk: req.Matk, Def: req.Def,
		Mdef: req.Mdef, Agi: req.Agi, Luk: req.Luk, Xp: req.Xp,
	})
	if err != nil {
		if errors.Is(err, domain.ErrCharacterNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "character not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to edit character"})
	}

	return c.JSON(http.StatusOK, presenter.EditCharacterRes{
		ID: editedCharacter.ID().String(), Name: editedCharacter.Name(), Job: editedCharacter.Job(),
		Hp: editedCharacter.Hp(), Atk: editedCharacter.Atk(), Matk: editedCharacter.Matk(),
		Def: editedCharacter.Def(), Mdef: editedCharacter.Mdef(), Agi: editedCharacter.Agi(),
		Luk: editedCharacter.Luk(), Xp: editedCharacter.Xp(), UserID: editedCharacter.UserID().String(),
		CreatedAt: editedCharacter.CreatedAt(),
	})
}
