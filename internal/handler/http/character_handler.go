package http

import (
	"backend/internal/domain"
	"backend/internal/handler/http/presenter"
	"backend/internal/usecase/character"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type CharacterHandler struct {
	createUc     *character.CreateCharacterUseCase
	statusEditUc *character.StatusEditCharacterUseCase
	jobChangeUc  *character.JobChangeCharacterUseCase
}

func NewCharacterHandler(createUc *character.CreateCharacterUseCase, statusEditUc *character.StatusEditCharacterUseCase, jobChangeUc *character.JobChangeCharacterUseCase) *CharacterHandler {
	return &CharacterHandler{
		createUc:     createUc,
		statusEditUc: statusEditUc,
		jobChangeUc:  jobChangeUc,
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

func (h *CharacterHandler) StatusEdit(c echo.Context) error {
	var req presenter.StatusEditCharacterReq
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

	editedCharacter, err := h.statusEditUc.Execute(c.Request().Context(), characterID, userID, req.Xp, character.StatusUpdate{
		Hp: req.Hp, Atk: req.Atk, Matk: req.Matk, Def: req.Def,
		Mdef: req.Mdef, Agi: req.Agi, Luk: req.Luk,
	})
	if err != nil {
		if errors.Is(err, domain.ErrCharacterNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "character not found"})
		}
		if errors.Is(err, domain.ErrInsufficientXP) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "insufficient xp"})
		}
		if errors.Is(err, domain.ErrInvalidXP) || errors.Is(err, domain.ErrInvalidStatus) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to edit character"})
	}

	return c.JSON(http.StatusOK, presenter.StatusEditCharacterRes{
		ID: editedCharacter.ID().String(), Name: editedCharacter.Name(), Job: editedCharacter.Job(),
		Hp: editedCharacter.Hp(), Atk: editedCharacter.Atk(), Matk: editedCharacter.Matk(),
		Def: editedCharacter.Def(), Mdef: editedCharacter.Mdef(), Agi: editedCharacter.Agi(),
		Luk: editedCharacter.Luk(), Xp: editedCharacter.Xp(), UserID: editedCharacter.UserID().String(),
		CreatedAt: editedCharacter.CreatedAt(),
	})
}

// JobChange godoc
// @Summary キャラクターを転職させる
// @Description ログイン中のユーザーが所有するキャラクターの職業を変更し、ステータスを初期化します。
// @Tags character
// @Accept json
// @Param request body presenter.JobChangeReq true "転職するキャラクターと職業"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /character/job [patch]
func (h *CharacterHandler) JobChange(c echo.Context) error {
	var req presenter.JobChangeReq
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

	job := strings.TrimSpace(req.Job)
	if job == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "job is req"})
	}

	err = h.jobChangeUc.Execute(c.Request().Context(), characterID, userID, job)
	if err != nil {
		if errors.Is(err, domain.ErrCharacterNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "character not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to edit character"})
	}

	return c.NoContent(http.StatusNoContent)
}
