package entity

import "time"

type UserAccount struct {
	ID           uint    `gorm:"primaryKey"`
	Username     string  `gorm:"uniqueIndex;not null"`
	PasswordHash string  `gorm:"not null"`
	Email        *string `gorm:"uniqueIndex"`
	PhoneNumber  *string
	FullName     string  `gorm:"not null"`
	IsActive     bool    `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (UserAccount) TableName() string {
	return "user_accounts"
}
