package user

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
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

func (s *Service) Login(ctx context.Context, req LoginRequest) (*User, error) {

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
