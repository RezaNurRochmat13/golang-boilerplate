package user

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, user *User) error {

	return r.db.WithContext(ctx).
		Create(user).
		Error
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {

	var user User

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
