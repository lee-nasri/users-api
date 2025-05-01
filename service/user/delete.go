package user

import (
	"context"
	"errors"

	"users-api/domain"
	"users-api/pkg/apperror"
	"users-api/pkg/logx"
)

func (s *Service) DeleteUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logx.Errorf(ctx, err, "user not found in database")
			return nil, apperror.NewErrUserNotFound()
		}
		logx.Errorf(ctx, err, "failed to get user by id")
		return nil, err
	}

	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		logx.Errorf(ctx, err, "Failed to delete user with ID %s", id)
		return nil, err
	}

	return user, nil
}
