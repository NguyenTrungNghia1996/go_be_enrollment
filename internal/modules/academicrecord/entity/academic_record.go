package entity

import "time"

type AcademicRecord struct {
	ID                  uint      `gorm:"primaryKey"`
	ApplicationID       uint      `gorm:"not null;uniqueIndex:idx_application_grade"`
	GradeLevel          int       `gorm:"not null;uniqueIndex:idx_application_grade"`
	SchoolName          string    `gorm:"size:255;not null"`
	SchoolLocation      string    `gorm:"size:255"`
	SchoolYear          string    `gorm:"size:50"`
	AcademicPerformance string    `gorm:"size:100"`
	ConductRating       string    `gorm:"size:100"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
