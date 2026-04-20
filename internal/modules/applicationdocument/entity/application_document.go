package entity

import "time"

type ApplicationDocument struct {
	ID            uint      `gorm:"primaryKey"`
	ApplicationID uint      `gorm:"not null"`
	DocumentType  string    `gorm:"size:100;not null"`
	FilePath      string    `gorm:"size:500;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
