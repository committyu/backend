package user

import (
	"context"

	"backend/internal/domain"
	"backend/internal/domain/repository"
)


type GetUserUsecase struct {
	userRepo repository.UserRepository
}

func NewGetUserUsecase(
	userRepo repository.UserRepository,
) *GetUserUsecase {
	return &GetUserUsecase{
		userRepo: userRepo,
	}
}

func (u *GetUserUsecase) Execute(
	ctx context.Context,
	userID domain.UserID,
) (*domain.User, error) {

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	
	return user, nil
}