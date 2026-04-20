package repository

import (
	"go_be_enrollment/internal/modules/application/dto"
	"go_be_enrollment/internal/modules/application/entity"
	
	admissionEntity "go_be_enrollment/internal/modules/admissionperiod/entity"
	provinceEntity "go_be_enrollment/internal/modules/province/entity"
	wardUnitEntity "go_be_enrollment/internal/modules/wardunit/entity"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	GetAdminList(filter *dto.ApplicationAdminFilter) ([]entity.Application, int64, error)
	GetAdminDetail(id uint) (*entity.Application, error)
	
	GetUserList(userID uint, page, limit int) ([]entity.Application, int64, error)
	GetUserDetail(id, userID uint) (*entity.Application, error)
	
	Create(app *entity.Application) error
	Update(app *entity.Application) error

	CheckAdmissionPeriodExists(id uint) bool
	CheckProvinceExists(id uint) bool
	CheckWardUnitExists(id, provinceID uint) bool
	CheckNationalIDExists(nationalID string, excludeID uint) bool
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) GetAdminList(filter *dto.ApplicationAdminFilter) ([]entity.Application, int64, error) {
	query := r.db.Model(&entity.Application{})

	if filter.AdmissionPeriodID != nil {
		query = query.Where("admission_period_id = ?", *filter.AdmissionPeriodID)
	}
	if filter.ApplicationStatus != "" {
		query = query.Where("application_status = ?", filter.ApplicationStatus)
	}
	if filter.IsPaid != nil {
		query = query.Where("is_paid = ?", *filter.IsPaid)
	}
	if filter.ProvinceID != nil {
		query = query.Where("province_id = ?", *filter.ProvinceID)
	}
	if filter.Keyword != "" {
		key := "%" + filter.Keyword + "%"
		query = query.Where("candidate_full_name LIKE ? OR national_id LIKE ? OR candidate_number LIKE ?", key, key, key)
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

	var list []entity.Application
	if err := query.Order("id desc").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *applicationRepository) GetAdminDetail(id uint) (*entity.Application, error) {
	var app entity.Application
	if err := r.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) GetUserList(userID uint, page, limit int) ([]entity.Application, int64, error) {
	query := r.db.Model(&entity.Application{}).Where("user_account_id = ?", userID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var list []entity.Application
	if err := query.Order("id desc").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *applicationRepository) GetUserDetail(id, userID uint) (*entity.Application, error) {
	var app entity.Application
	if err := r.db.Where("id = ? AND user_account_id = ?", id, userID).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) Create(app *entity.Application) error {
	return r.db.Create(app).Error
}

func (r *applicationRepository) Update(app *entity.Application) error {
	return r.db.Save(app).Error
}

func (r *applicationRepository) CheckAdmissionPeriodExists(id uint) bool {
	var count int64
	r.db.Model(&admissionEntity.AdmissionPeriod{}).Where("id = ?", id).Count(&count)
	return count > 0
}

func (r *applicationRepository) CheckProvinceExists(id uint) bool {
	var count int64
	r.db.Model(&provinceEntity.Province{}).Where("id = ?", id).Count(&count)
	return count > 0
}

func (r *applicationRepository) CheckWardUnitExists(id, provinceID uint) bool {
	var count int64
	r.db.Model(&wardUnitEntity.WardUnit{}).Where("id = ? AND province_id = ?", id, provinceID).Count(&count)
	return count > 0
}

func (r *applicationRepository) CheckNationalIDExists(nationalID string, excludeID uint) bool {
	var count int64
	q := r.db.Model(&entity.Application{}).Where("national_id = ?", nationalID)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	q.Count(&count)
	return count > 0
}
