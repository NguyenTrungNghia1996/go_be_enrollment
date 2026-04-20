package entity

import "time"

type ExamRoom struct {
	ID        uint      `gorm:"primaryKey"`
	RoomName  string    `gorm:"size:255;not null"`
	Location  string    `gorm:"size:255;not null"`
	Capacity  int       `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
