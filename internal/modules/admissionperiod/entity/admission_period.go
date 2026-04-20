package entity

import "time"

type AdmissionPeriod struct {
	ID         uint       `gorm:"primaryKey"`
	Name       string     `gorm:"size:255;not null"`
	SchoolYear string     `gorm:"size:50;not null"`
	StartDate  *time.Time `gorm:"type:date"`
	EndDate    *time.Time `gorm:"type:date"`
	ExamFee    float64    `gorm:"type:decimal(15,2);default:0"`
	IsOpen     bool       `gorm:"default:false"`
	Notes      string     `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
