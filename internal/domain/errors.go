package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrCharacterNotFound = errors.New("character not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidToken      = errors.New("invalid token")
)
