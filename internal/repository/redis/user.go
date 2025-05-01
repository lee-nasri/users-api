package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"users-api/domain"
	"users-api/pkg/helper"
)

type UserRepository struct {
	wc   *WrapperClient
	rCfg UserConfig
}

type UserConfig struct {
	KeyPrefix    string
	Index        string
	DefaultLimit int
}

func NewUserRepository(wc *WrapperClient, rCfg *UserConfig) *UserRepository {
	return &UserRepository{
		wc:   wc,
		rCfg: *rCfg,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, data domain.User) (*domain.User, error) {
	data.CreatedAt = time.Now().UnixMilli()
	rKey := fmt.Sprintf("%s:%s", r.rCfg.KeyPrefix, data.ID)

	// Marshal
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// save to redis
	err = r.wc.Client.JSONSet(ctx, rKey, ".", jsonData).Err()
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *UserRepository) GetUserByUserID(ctx context.Context, id string) (*domain.User, error) {
	rKey := fmt.Sprintf("%s:%s", r.rCfg.KeyPrefix, id)
	var user domain.User

	// get from redis
	res, err := r.wc.Client.JSONGet(ctx, rKey, ".").Result()
	if err != nil {
		return nil, err
	}

	if res == "" {
		return nil, domain.ErrUserNotFound
	}

	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserBySurname(ctx context.Context, surname string) (*domain.User, error) {
	// get from redis
	_, values, err := r.wc.SearchJSONDataText(ctx, r.rCfg.Index, "surname", surname)
	if err != nil {
		return nil, err
	}

	var users = make([]domain.User, 0, len(values))

	for _, v := range values {
		var d domain.User

		if err = json.Unmarshal([]byte(v), &d); err != nil {
			return nil, err
		}
		users = append(users, d)
	}

	if len(users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	// If more than one symbol is found, return the first one. Note that the symbol_name is forced to be unique when created through the API.
	return &users[0], nil
}

func (r *UserRepository) GetUserByLastname(ctx context.Context, lastname string) (*domain.User, error) {
	// get from redis
	_, values, err := r.wc.SearchJSONDataText(ctx, r.rCfg.Index, "lastname", lastname)
	if err != nil {
		return nil, err
	}

	var users = make([]domain.User, 0, len(values))

	for _, v := range values {
		var d domain.User

		if err = json.Unmarshal([]byte(v), &d); err != nil {
			return nil, err
		}
		users = append(users, d)
	}

	if len(users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	// If more than one symbol is found, return the first one. Note that the symbol_name is forced to be unique when created through the API.
	return &users[0], nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, data domain.User) (*domain.User, error) {
	data.UpdatedAt = helper.ToPtr(time.Now().UnixMilli())
	rKey := fmt.Sprintf("%s:%s", r.rCfg.KeyPrefix, id)

	// Marshal
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// save to redis
	err = r.wc.Client.JSONSet(ctx, rKey, ".", jsonData).Err()
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	rKey := fmt.Sprintf("%s:%s", r.rCfg.KeyPrefix, id)

	// delete from redis
	err := r.wc.Client.Del(ctx, rKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) IsUserIDAlreadyExists(ctx context.Context, id string) (bool, error) {
	keys, _, err := r.wc.SearchJSONDataText(ctx, r.rCfg.Index, "id", id)
	if err != nil {
		return false, err
	}

	if len(keys) > 0 {
		return true, nil
	}

	return false, nil
}
