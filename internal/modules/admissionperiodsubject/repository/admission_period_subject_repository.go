package repository

import (
	"go_be_enrollment/internal/modules/admissionperiodsubject/entity"

	"gorm.io/gorm"
)

type AdmissionPeriodSubjectRepository interface {
	GetByAdmissionPeriodID(periodID uint) ([]entity.AdmissionPeriodSubject, error)
	ReplaceSubjects(periodID uint, subjects []entity.AdmissionPeriodSubject) error
}

type admissionPeriodSubjectRepository struct {
	db *gorm.DB
}

func NewAdmissionPeriodSubjectRepository(db *gorm.DB) AdmissionPeriodSubjectRepository {
	return &admissionPeriodSubjectRepository{db: db}
}

func (r *admissionPeriodSubjectRepository) GetByAdmissionPeriodID(periodID uint) ([]entity.AdmissionPeriodSubject, error) {
	var list []entity.AdmissionPeriodSubject
	if err := r.db.Preload("Subject").Where("admission_period_id = ?", periodID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *admissionPeriodSubjectRepository) ReplaceSubjects(periodID uint, subjects []entity.AdmissionPeriodSubject) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing
		if err := tx.Where("admission_period_id = ?", periodID).Delete(&entity.AdmissionPeriodSubject{}).Error; err != nil {
			return err
		}

		// Insert new
		if len(subjects) > 0 {
			if err := tx.Create(&subjects).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
