package tests

import (
	"context"
	"errors"
	"testing"

	"golang-boilerplate-example/module/user"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	createFn      func(ctx context.Context, user *user.User) error
	findByEmailFn func(ctx context.Context, email string) (*user.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, u *user.User) error {
	return m.createFn(ctx, u)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return m.findByEmailFn(ctx, email)
}

func TestService_Register_Success(t *testing.T) {

	repo := &mockUserRepository{
		createFn: func(ctx context.Context, u *user.User) error {

			if u.Email != "test@mail.com" {
				t.Fatalf("unexpected email: %s", u.Email)
			}

			if u.Name != "Reja" {
				t.Fatalf("unexpected name: %s", u.Name)
			}

			if u.Password == "" {
				t.Fatal("password should be hashed")
			}

			return nil
		},
	}

	service := user.NewService(repo)

	err := service.Register(context.Background(), user.RegisterRequest{
		Email:    "test@mail.com",
		Name:     "Reja",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_Register_RepositoryError(t *testing.T) {

	repoErr := errors.New("db error")

	repo := &mockUserRepository{
		createFn: func(ctx context.Context, u *user.User) error {
			return repoErr
		},
	}

	service := user.NewService(repo)

	err := service.Register(context.Background(), user.RegisterRequest{
		Email:    "test@mail.com",
		Name:     "Reja",
		Password: "password123",
	})

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestService_Login_Success(t *testing.T) {

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	expectedUser := &user.User{
		ID:       1,
		Email:    "test@mail.com",
		Name:     "Reja",
		Password: string(hash),
	}

	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*user.User, error) {

			if email != "test@mail.com" {
				t.Fatalf("unexpected email: %s", email)
			}

			return expectedUser, nil
		},
	}

	service := user.NewService(repo)

	u, err := service.Login(context.Background(), user.LoginRequest{
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Email != expectedUser.Email {
		t.Fatalf("expected %s got %s", expectedUser.Email, u.Email)
	}
}

func TestService_Login_UserNotFound(t *testing.T) {

	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
			return nil, errors.New("not found")
		},
	}

	service := user.NewService(repo)

	_, err := service.Login(context.Background(), user.LoginRequest{
		Email:    "test@mail.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials got %v", err)
	}
}

func TestService_Login_InvalidPassword(t *testing.T) {

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

	existingUser := &user.User{
		ID:       1,
		Email:    "test@mail.com",
		Name:     "Reja",
		Password: string(hash),
	}

	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
			return existingUser, nil
		},
	}

	service := user.NewService(repo)

	_, err := service.Login(context.Background(), user.LoginRequest{
		Email:    "test@mail.com",
		Password: "wrong-password",
	})

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials got %v", err)
	}
}
