package auth

import (
	"backend/internal/domain"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct{}

var _ domain.JWTService = (*JWTService)(nil)

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (j *JWTService) Generate(userID string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	return token.SignedString([]byte(secret))
}