package repository

import (
	"errors"

	"go_be_enrollment/internal/modules/academicrecord/entity"
	appEntity "go_be_enrollment/internal/modules/application/entity"

	"gorm.io/gorm"
)

type AcademicRecordRepository interface {
	GetByApplicationID(appID uint) ([]entity.AcademicRecord, error)
	VerifyApplicationAccess(appID, userID uint) (bool, string, error) // return isOwner, applicationStatus, error
	FindByID(id uint) (*entity.AcademicRecord, error)
	Create(record *entity.AcademicRecord) error
	Update(record *entity.AcademicRecord) error
	Delete(id uint) error
	CheckDuplicateGrade(appID uint, gradeLevel int, excludeID uint) bool
}

type academicRecordRepository struct {
	db *gorm.DB
}

func NewAcademicRecordRepository(db *gorm.DB) AcademicRecordRepository {
	return &academicRecordRepository{db: db}
}

func (r *academicRecordRepository) GetByApplicationID(appID uint) ([]entity.AcademicRecord, error) {
	var list []entity.AcademicRecord
	if err := r.db.Where("application_id = ?", appID).Order("grade_level ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *academicRecordRepository) VerifyApplicationAccess(appID, userID uint) (bool, string, error) {
	var app appEntity.Application
	if err := r.db.Select("user_account_id", "application_status").Where("id = ?", appID).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", nil
		}
		return false, "", err
	}
	
	isOwner := app.UserAccountID == userID
	return isOwner, app.ApplicationStatus, nil
}

func (r *academicRecordRepository) FindByID(id uint) (*entity.AcademicRecord, error) {
	var record entity.AcademicRecord
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *academicRecordRepository) Create(record *entity.AcademicRecord) error {
	return r.db.Create(record).Error
}

func (r *academicRecordRepository) Update(record *entity.AcademicRecord) error {
	return r.db.Save(record).Error
}

func (r *academicRecordRepository) Delete(id uint) error {
	return r.db.Delete(&entity.AcademicRecord{}, id).Error
}

func (r *academicRecordRepository) CheckDuplicateGrade(appID uint, gradeLevel int, excludeID uint) bool {
	var count int64
	q := r.db.Model(&entity.AcademicRecord{}).Where("application_id = ? AND grade_level = ?", appID, gradeLevel)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	q.Count(&count)
	return count > 0
}
