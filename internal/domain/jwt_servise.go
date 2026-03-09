package domain

type JWTService interface {
	Generate(userID string) (string, error)
}