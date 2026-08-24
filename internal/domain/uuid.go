package domain

import (
	"github.com/google/uuid"
)

type UserID string

type CharacterID string

func newID[T ~string]() T {
	return T(uuid.NewString())
}

func parseID[T ~string](s string) (T, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		var zero T
		return zero, err
	}

	return T(id.String()), nil
}

func isValidID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func NewUserID() UserID {
	return newID[UserID]()
}

func (id UserID) String() string {
	return string(id)
}

func ParseUserID(s string) (UserID, error) {
	return parseID[UserID](s)
}

func IsValidUserID(id string) bool {
	return isValidID(id)
}

func NewCharacterID() CharacterID {
	return newID[CharacterID]()
}

func (id CharacterID) String() string {
	return string(id)
}

func ParseCharacterID(s string) (CharacterID, error) {
	return parseID[CharacterID](s)
}

func IsValidCharacterID(id string) bool {
	return isValidID(id)
}
