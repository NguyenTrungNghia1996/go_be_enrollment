package repository

import (
	authentity "go_be_enrollment/internal/modules/auth/entity"
	"go_be_enrollment/internal/modules/useraccount/dto"

	"gorm.io/gorm"
)

type UserAccountRepository interface {
	GetList(filter *dto.UserAccountFilter) ([]authentity.UserAccount, int64, error)
	FindByID(id uint) (*authentity.UserAccount, error)
	Update(user *authentity.UserAccount) error
}

type userAccountRepository struct {
	db *gorm.DB
}

func NewUserAccountRepository(db *gorm.DB) UserAccountRepository {
	return &userAccountRepository{db: db}
}

func (r *userAccountRepository) GetList(filter *dto.UserAccountFilter) ([]authentity.UserAccount, int64, error) {
	query := r.db.Model(&authentity.UserAccount{})

	if filter.Keyword != "" {
		key := "%" + filter.Keyword + "%"
		query = query.Where("username LIKE ? OR full_name LIKE ? OR email LIKE ? OR phone_number LIKE ?", key, key, key, key)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	switch {
	case limit <= 0:
		limit = 10
	case limit > 100:
		limit = 100
	}
	offset := (page - 1) * limit

	var list []authentity.UserAccount
	if err := query.Order("id desc").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *userAccountRepository) FindByID(id uint) (*authentity.UserAccount, error) {
	var user authentity.UserAccount
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userAccountRepository) Update(user *authentity.UserAccount) error {
	return r.db.Save(user).Error
}
