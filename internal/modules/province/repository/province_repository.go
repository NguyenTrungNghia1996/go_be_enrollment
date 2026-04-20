package repository

import (
	"go_be_enrollment/internal/modules/province/dto"
	"go_be_enrollment/internal/modules/province/entity"

	"gorm.io/gorm"
)

type ProvinceRepository interface {
	GetList(filter *dto.ProvinceFilter) ([]entity.Province, int64, error)
	FindByID(id uint) (*entity.Province, error)
	Create(p *entity.Province) error
	Update(p *entity.Province) error
	CheckCodeExists(code string, excludeID uint) bool
	CheckNameExists(name string, excludeID uint) bool
	GetActiveList() ([]entity.Province, error)
}

type provinceRepository struct {
	db *gorm.DB
}

func NewProvinceRepository(db *gorm.DB) ProvinceRepository {
	return &provinceRepository{db: db}
}

func (r *provinceRepository) GetList(filter *dto.ProvinceFilter) ([]entity.Province, int64, error) {
	query := r.db.Model(&entity.Province{})

	if filter.Keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
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
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}
	offset := (page - 1) * limit

	var provinces []entity.Province
	if err := query.Order("id desc").Offset(offset).Limit(limit).Find(&provinces).Error; err != nil {
		return nil, 0, err
	}

	return provinces, total, nil
}

func (r *provinceRepository) FindByID(id uint) (*entity.Province, error) {
	var province entity.Province
	if err := r.db.First(&province, id).Error; err != nil {
		return nil, err
	}
	return &province, nil
}

func (r *provinceRepository) Create(p *entity.Province) error {
	return r.db.Create(p).Error
}

func (r *provinceRepository) Update(p *entity.Province) error {
	return r.db.Save(p).Error
}

func (r *provinceRepository) CheckCodeExists(code string, excludeID uint) bool {
	var count int64
	r.db.Model(&entity.Province{}).Where("code = ? AND id != ?", code, excludeID).Count(&count)
	return count > 0
}

func (r *provinceRepository) CheckNameExists(name string, excludeID uint) bool {
	var count int64
	r.db.Model(&entity.Province{}).Where("name = ? AND id != ?", name, excludeID).Count(&count)
	return count > 0
}

func (r *provinceRepository) GetActiveList() ([]entity.Province, error) {
	var provinces []entity.Province
	if err := r.db.Where("is_active = ?", true).Order("name ASC").Find(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}
