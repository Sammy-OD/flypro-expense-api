package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"flypro/internal/models"
	"flypro/internal/repository"

	"github.com/redis/go-redis/v9"
)

type UserService interface {
	CreateUser(ctx context.Context, email, name string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
}

type userService struct {
	repo  repository.UserRepository
	cache *redis.Client
}

func NewUserService(repo repository.UserRepository, cache *redis.Client) UserService {
	return &userService{repo: repo, cache: cache}
}

func (s *userService) CreateUser(ctx context.Context, email, name string) (*models.User, error) {
	u := &models.User{Email: email, Name: name}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	// Invalidate cache
	_ = s.cache.Del(ctx, s.userCacheKey(u.ID)).Err()
	return u, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	key := s.userCacheKey(id)
	if b, err := s.cache.Get(ctx, key).Bytes(); err == nil {
		var u models.User
		if json.Unmarshal(b, &u) == nil {
			return &u, nil
		}
	}
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if b, _ := json.Marshal(u); b != nil {
		_ = s.cache.Set(ctx, key, b, time.Hour).Err()
	}
	return u, nil
}

func (s *userService) userCacheKey(id uint) string { return fmt.Sprintf("user:%d", id) }
