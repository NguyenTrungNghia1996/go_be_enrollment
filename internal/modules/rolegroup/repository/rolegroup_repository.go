package repository

import (
	"strings"

	"go_be_enrollment/internal/modules/adminauth/entity"
	"go_be_enrollment/internal/modules/rolegroup/dto"

	"gorm.io/gorm"
)

type RoleGroupRepository interface {
	FindAll(filter *dto.RoleGroupFilter) ([]entity.RoleGroup, int64, error)
	FindByID(id uint) (*entity.RoleGroup, error)
	Create(role *entity.RoleGroup) error
	Update(role *entity.RoleGroup) error
	CheckCodeExists(code string, excludeID uint) bool
	CheckIsAssigned(id uint) bool
}

type roleGroupRepository struct {
	db *gorm.DB
}

func NewRoleGroupRepository(db *gorm.DB) RoleGroupRepository {
	return &roleGroupRepository{db: db}
}

func (r *roleGroupRepository) FindAll(filter *dto.RoleGroupFilter) ([]entity.RoleGroup, int64, error) {
	var roles []entity.RoleGroup
	var total int64

	query := r.db.Model(&entity.RoleGroup{})

	if filter.Keyword != "" {
		kw := "%" + strings.ToLower(filter.Keyword) + "%"
		query = query.Where("LOWER(code) LIKE ? OR LOWER(name) LIKE ?", kw, kw)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	err = query.Order("id desc").Offset(offset).Limit(limit).Find(&roles).Error
	return roles, total, err
}

func (r *roleGroupRepository) FindByID(id uint) (*entity.RoleGroup, error) {
	var role entity.RoleGroup
	err := r.db.First(&role, id).Error
	return &role, err
}

func (r *roleGroupRepository) Create(role *entity.RoleGroup) error {
	return r.db.Create(role).Error
}

func (r *roleGroupRepository) Update(role *entity.RoleGroup) error {
	return r.db.Save(role).Error
}

func (r *roleGroupRepository) CheckCodeExists(code string, excludeID uint) bool {
	var count int64
	r.db.Model(&entity.RoleGroup{}).Where("code = ? AND id != ?", code, excludeID).Count(&count)
	return count > 0
}

func (r *roleGroupRepository) CheckIsAssigned(id uint) bool {
	var count int64
	r.db.Model(&entity.AdminUserRoleGroup{}).Where("role_group_id = ?", id).Count(&count)
	return count > 0
}
