package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrCharacterNotFound = errors.New("character not found")
	ErrInsufficientXP    = errors.New("insufficient xp")
	ErrInvalidXP         = errors.New("xp must be greater than zero")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidToken      = errors.New("invalid token")
)
