package repository

import (
	"go_be_enrollment/internal/modules/admissionperiod/dto"
	"go_be_enrollment/internal/modules/admissionperiod/entity"

	"gorm.io/gorm"
)

type AdmissionPeriodRepository interface {
	GetList(filter *dto.AdmissionPeriodFilter) ([]entity.AdmissionPeriod, int64, error)
	FindByID(id uint) (*entity.AdmissionPeriod, error)
	Create(ap *entity.AdmissionPeriod) error
	Update(ap *entity.AdmissionPeriod) error
	GetOpenPeriods() ([]entity.AdmissionPeriod, error)
}

type admissionPeriodRepository struct {
	db *gorm.DB
}

func NewAdmissionPeriodRepository(db *gorm.DB) AdmissionPeriodRepository {
	return &admissionPeriodRepository{db: db}
}

func (r *admissionPeriodRepository) GetList(filter *dto.AdmissionPeriodFilter) ([]entity.AdmissionPeriod, int64, error) {
	query := r.db.Model(&entity.AdmissionPeriod{})

	if filter.Keyword != "" {
		query = query.Where("name LIKE ? OR notes LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	if filter.SchoolYear != "" {
		query = query.Where("school_year = ?", filter.SchoolYear)
	}

	if filter.IsOpen != nil {
		query = query.Where("is_open = ?", *filter.IsOpen)
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

	var list []entity.AdmissionPeriod
	if err := query.Order("id desc").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *admissionPeriodRepository) FindByID(id uint) (*entity.AdmissionPeriod, error) {
	var ap entity.AdmissionPeriod
	if err := r.db.First(&ap, id).Error; err != nil {
		return nil, err
	}
	return &ap, nil
}

func (r *admissionPeriodRepository) Create(ap *entity.AdmissionPeriod) error {
	return r.db.Create(ap).Error
}

func (r *admissionPeriodRepository) Update(ap *entity.AdmissionPeriod) error {
	return r.db.Save(ap).Error
}

func (r *admissionPeriodRepository) GetOpenPeriods() ([]entity.AdmissionPeriod, error) {
	var list []entity.AdmissionPeriod
	if err := r.db.Where("is_open = ?", true).Order("start_date ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
