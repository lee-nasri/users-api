package user

import (
	"context"
	"errors"

	"users-api/domain"
	"users-api/pkg/logx"
)

func (s *Service) UpdateUser(ctx context.Context, id string, data domain.UpdateUserRequest) (*domain.User, error) {
	logx.Infof(ctx, "Update user with ID: %s, data: %+v", id, data)

	userData, err := s.userRepo.GetUserByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			logx.Errorf(ctx, err, "[!] User not found")
			return nil, domain.ErrUserNotFound
		}
		logx.Errorf(ctx, err, "[!] Error fetching user")
		return nil, err
	}

	userData = parseValue(*userData, data)

	updatedUser, err := s.userRepo.UpdateUser(ctx, id, *userData)
	if err != nil {
		logx.Errorf(ctx, err, "[!] Error updating user")
		return nil, err
	}

	return updatedUser, nil
}

func parseValue(userData domain.User, input domain.UpdateUserRequest) *domain.User {
	if input.Surname != nil {
		userData.Surname = *input.Surname
	}
	if input.Lastname != nil {
		userData.Lastname = *input.Lastname
	}
	if input.Age != nil {
		userData.Age = *input.Age
	}
	if input.Email != nil {
		userData.Email = *input.Email
	}
	if input.Phone != nil {
		userData.Phone = *input.Phone
	}

	return &userData
}
