package entity

import "time"

type Province struct {
	ID        uint      `gorm:"primaryKey"`
	Code      string    `gorm:"size:50;uniqueIndex;not null"`
	Name      string    `gorm:"size:255;uniqueIndex;not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
