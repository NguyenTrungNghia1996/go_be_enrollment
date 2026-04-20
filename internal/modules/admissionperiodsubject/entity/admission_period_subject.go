package entity

import (
	"time"
	"go_be_enrollment/internal/modules/subject/entity"
)

type AdmissionPeriodSubject struct {
	ID                uint           `gorm:"primaryKey"`
	AdmissionPeriodID uint           `gorm:"not null;uniqueIndex:idx_admission_period_subject"`
	SubjectID         uint           `gorm:"not null;uniqueIndex:idx_admission_period_subject"`
	Weight            float64        `gorm:"type:decimal(5,2);default:1.0"`
	IsRequired        bool           `gorm:"default:true"`
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Relations
	Subject           *entity.Subject `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
