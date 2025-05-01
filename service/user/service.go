package user

import (
	"context"

	"users-api/domain"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, data domain.User) (*domain.User, error)
	GetUserByUserID(ctx context.Context, id string) (*domain.User, error)
	GetUserBySurname(ctx context.Context, surname string) (*domain.User, error)
	GetUserByLastname(ctx context.Context, lastname string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, data domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error

	IsUserIDAlreadyExists(ctx context.Context, id string) (bool, error)
}

type Service struct {
	userRepo IUserRepository
}

func NewService(userRepo IUserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}
