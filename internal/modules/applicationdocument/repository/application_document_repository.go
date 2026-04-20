package repository

import (
	"errors"

	"go_be_enrollment/internal/modules/applicationdocument/entity"
	appEntity "go_be_enrollment/internal/modules/application/entity"

	"gorm.io/gorm"
)

type ApplicationDocumentRepository interface {
	GetByApplicationID(appID uint) ([]entity.ApplicationDocument, error)
	VerifyApplicationAccess(appID, userID uint) (bool, string, error) // isOwner, status, err
	FindByID(id uint) (*entity.ApplicationDocument, error)
	Create(doc *entity.ApplicationDocument) error
	Delete(id uint) error
}

type applicationDocumentRepository struct {
	db *gorm.DB
}

func NewApplicationDocumentRepository(db *gorm.DB) ApplicationDocumentRepository {
	return &applicationDocumentRepository{db: db}
}

func (r *applicationDocumentRepository) GetByApplicationID(appID uint) ([]entity.ApplicationDocument, error) {
	var list []entity.ApplicationDocument
	if err := r.db.Where("application_id = ?", appID).Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *applicationDocumentRepository) VerifyApplicationAccess(appID, userID uint) (bool, string, error) {
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

func (r *applicationDocumentRepository) FindByID(id uint) (*entity.ApplicationDocument, error) {
	var doc entity.ApplicationDocument
	if err := r.db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *applicationDocumentRepository) Create(doc *entity.ApplicationDocument) error {
	return r.db.Create(doc).Error
}

func (r *applicationDocumentRepository) Delete(id uint) error {
	return r.db.Delete(&entity.ApplicationDocument{}, id).Error
}
