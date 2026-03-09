package auth

import "backend/internal/domain"

type GenerateTokenUsecase struct {
	jwt domain.JWTService
}

func NewGenerateTokenUsecase(jwt domain.JWTService) *GenerateTokenUsecase {
	return &GenerateTokenUsecase{
		jwt: jwt,
	}
}

func (u *GenerateTokenUsecase) Execute(userID domain.UserID) (string, error) {
	return u.jwt.Generate(userID.String())
}