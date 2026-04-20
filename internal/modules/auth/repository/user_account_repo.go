package repository

import (
	"context"

	"go_be_enrollment/internal/modules/auth/entity"
	"gorm.io/gorm"
)

type UserAccountRepository interface {
	Create(ctx context.Context, user *entity.UserAccount) error
	FindByUsername(ctx context.Context, username string) (*entity.UserAccount, error)
	FindByEmail(ctx context.Context, email string) (*entity.UserAccount, error)
	FindByID(ctx context.Context, id uint) (*entity.UserAccount, error)
}

type userAccountRepo struct {
	db *gorm.DB
}

func NewUserAccountRepository(db *gorm.DB) UserAccountRepository {
	return &userAccountRepo{db: db}
}

func (r *userAccountRepo) Create(ctx context.Context, user *entity.UserAccount) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userAccountRepo) FindByUsername(ctx context.Context, username string) (*entity.UserAccount, error) {
	var user entity.UserAccount
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userAccountRepo) FindByEmail(ctx context.Context, email string) (*entity.UserAccount, error) {
	var user entity.UserAccount
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userAccountRepo) FindByID(ctx context.Context, id uint) (*entity.UserAccount, error) {
	var user entity.UserAccount
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
