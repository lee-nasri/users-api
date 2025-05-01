package user

import (
	"context"
	"errors"
	"log"

	"users-api/domain"
	_const "users-api/pkg/const"
	"users-api/pkg/helper"
	"users-api/pkg/logx"
)

func (s *Service) CreateUser(ctx context.Context, input domain.CreateUserRequest) (*domain.User, error) {
	logx.Infof(ctx, "Create new user input: %+v", input)
	var (
		retryCount int
		userID     string
		err        error
	)

	for {
		userID, err = helper.GenerateID()
		if err != nil {
			logx.Errorf(ctx, err, "[!] Error generating user ID")
			return nil, err
		}

		isUserIDExist, err := s.userRepo.IsUserIDAlreadyExists(ctx, userID)
		if err != nil {
			logx.Errorf(ctx, err, "[! Error checking user ID existence")
			return nil, err
		}

		if !isUserIDExist {
			break
		}

		log.Printf("User ID %s already exists, generating a new one", userID)
		if retryCount >= _const.RetryGenUserIDLimit {
			e := errors.New("failed to generate unique user ID, retry limit reached")
			logx.Errorf(ctx, e, "[!] Error generating unique user ID, retry limit reached")
			return nil, e
		}
		retryCount++
	}

	res, err := s.userRepo.CreateUser(ctx, domain.User{
		ID:       userID,
		Surname:  input.Surname,
		Lastname: input.Lastname,
		Age:      input.Age,
		Email:    input.Email,
		Phone:    input.Phone,
	})

	if err != nil {
		logx.Errorf(ctx, err, "Error creating user")
		return nil, err
	}

	return res, nil
}
