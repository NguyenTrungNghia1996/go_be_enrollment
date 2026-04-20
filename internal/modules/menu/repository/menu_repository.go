package repository

import (
	"go_be_enrollment/internal/modules/adminauth/entity"
	"gorm.io/gorm"
)

type MenuRepository interface {
	FindAll() ([]entity.Menu, error)
	FindByID(id uint) (*entity.Menu, error)
	Create(m *entity.Menu) error
	Update(m *entity.Menu) error
	Delete(m *entity.Menu) error
	CheckMenuKeyExists(key string, excludeID uint) bool
	CountChildren(parentID uint) (int64, error)
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) FindAll() ([]entity.Menu, error) {
	var menus []entity.Menu
	// sort by sort_order ASC, then ID ASC
	err := r.db.Order("sort_order ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByID(id uint) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.db.First(&menu, id).Error
	return &menu, err
}

func (r *menuRepository) Create(m *entity.Menu) error {
	return r.db.Create(m).Error
}

func (r *menuRepository) Update(m *entity.Menu) error {
	return r.db.Save(m).Error
}

func (r *menuRepository) Delete(m *entity.Menu) error {
	return r.db.Delete(m).Error
}

func (r *menuRepository) CheckMenuKeyExists(key string, excludeID uint) bool {
	var count int64
	r.db.Model(&entity.Menu{}).Where("menu_key = ? AND id != ?", key, excludeID).Count(&count)
	return count > 0
}

func (r *menuRepository) CountChildren(parentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Menu{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}
