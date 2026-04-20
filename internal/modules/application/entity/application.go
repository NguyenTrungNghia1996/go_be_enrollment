package entity

import "time"

type Application struct {
	ID                 uint       `gorm:"primaryKey"`
	UserAccountID      uint       `gorm:"not null"`
	AdmissionPeriodID  uint       `gorm:"not null;uniqueIndex:idx_admission_candidate"`
	CandidateFullName  string     `gorm:"size:255;not null"`
	DateOfBirth        *time.Time `gorm:"type:date;not null"`
	Gender             string     `gorm:"type:enum('Male','Female','Other');not null"`
	PlaceOfBirth       string     `gorm:"size:255"`
	Ethnicity          string     `gorm:"size:100"`
	NationalID         string     `gorm:"size:50;uniqueIndex;not null"`
	ProvinceID         *uint      
	WardUnitID         *uint      
	AddressLine        string     `gorm:"size:255"`
	ContactFullName    string     `gorm:"size:255;not null"`
	ContactPhoneNumber string     `gorm:"size:20;not null"`
	ApplicationStatus  string     `gorm:"type:enum('Draft','Submitted','Approved','Rejected');default:'Draft'"`
	IsPaid             bool       `gorm:"default:false"`
	SubmittedAt        *time.Time 
	CandidateNumber    *string    `gorm:"size:50;uniqueIndex:idx_admission_candidate"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
