package Usecases

import (
	"context"
	"errors"
	"taskmanager/Domain"
	"taskmanager/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserUsecase defines the use case interface for user operations
type UserUsecase interface {
	RegisterUser(ctx context.Context, user Domain.User) error
	AuthenticateUser(ctx context.Context, username, password string) (*Domain.User, error)
	GetUserByID(ctx context.Context, id string) (*Domain.User, error)
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo Repositories.UserRepository
}

// NewUserUsecase creates a new UserUsecase
func NewUserUsecase(userRepo Repositories.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) RegisterUser(ctx context.Context, user Domain.User) error {
	return u.userRepo.RegisterUser(ctx, user)
}

func (u *userUsecase) AuthenticateUser(ctx context.Context, username, password string) (*Domain.User, error) {
	return u.userRepo.AuthenticateUser(ctx, username, password)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*Domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return u.userRepo.GetUserByID(ctx, objID)
}
