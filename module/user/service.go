package user

import (
	"context"
	"errors"
	"fmt"
	"golang-boilerplate-example/internal/auth"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo  Repository
	redis *redis.Client
}

func NewService(repo Repository, redis *redis.Client) *Service {
	return &Service{repo, redis}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hash),
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(uint(user.ID))
	if err != nil {
		return "", err
	}

	// store session in redis
	key := fmt.Sprintf("auth:session:%d", user.ID)

	err = s.redis.Set(ctx, key, token, time.Hour*24).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}
