package domain

import (
	"github.com/google/uuid"
)

type UserID string

func NewUserID() UserID {
	return UserID(uuid.NewString())
}

func (id UserID) String() string {
    return string(id)
}

func UserIDFromString(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}
	return UserID(id.String()), nil
}

func IsValidUserID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}