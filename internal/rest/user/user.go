package user

import (
	"context"

	"users-api/domain"
	"users-api/pkg/apperror"
)

type IUserService interface {
	CreateUser(ctx context.Context, input domain.CreateUserRequest) (*domain.User, error)
	GetUserByUserID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, data domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) (*domain.User, error)
}

type IValidator interface {
	ValidateStruct(v interface{}) error
}

type Handler struct {
	userSvc   IUserService
	validator IValidator
}

func NewHandler(userSvc IUserService, validator IValidator) *Handler {
	return &Handler{
		userSvc:   userSvc,
		validator: validator,
	}
}

func (s *Handler) validateBodyParser(body interface{}) error {
	err := s.validator.ValidateStruct(body)
	if err != nil {
		_, ok := apperror.IsAppError(err)
		if ok {
			return err
		}
		return apperror.NewInvalidRequestFromErr(err)
	}
	return nil
}
