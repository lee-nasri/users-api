package user

import (
	"context"
	"errors"

	"users-api/domain"
	"users-api/pkg/apperror"
	"users-api/pkg/logx"
)

func (s *Service) GetUserByUserID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logx.Errorf(ctx, err, "user not found in database")
			return nil, apperror.NewErrUserNotFound()
		}
		logx.Errorf(ctx, err, "failed to get user by id")
		return nil, err
	}

	return user, nil
}
