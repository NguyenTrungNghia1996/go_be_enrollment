package dto

type ApplicationReq struct {
	AdmissionPeriodID  uint    `json:"admission_period_id" validate:"required"`
	CandidateFullName  string  `json:"candidate_full_name" validate:"required"`
	DateOfBirth        string  `json:"date_of_birth" validate:"required"` // YYYY-MM-DD
	Gender             string  `json:"gender" validate:"required,oneof=Male Female Other"`
	PlaceOfBirth       string  `json:"place_of_birth"`
	Ethnicity          string  `json:"ethnicity"`
	NationalID         string  `json:"national_id" validate:"required"`
	ProvinceID         *uint   `json:"province_id"`
	WardUnitID         *uint   `json:"ward_unit_id"`
	AddressLine        string  `json:"address_line"`
	ContactFullName    string  `json:"contact_full_name" validate:"required"`
	ContactPhoneNumber string  `json:"contact_phone_number" validate:"required"`
}

type ApplicationRes struct {
	ID                 uint    `json:"id"`
	UserAccountID      uint    `json:"user_account_id"`
	AdmissionPeriodID  uint    `json:"admission_period_id"`
	CandidateFullName  string  `json:"candidate_full_name"`
	DateOfBirth        string  `json:"date_of_birth"`
	Gender             string  `json:"gender"`
	PlaceOfBirth       string  `json:"place_of_birth"`
	Ethnicity          string  `json:"ethnicity"`
	NationalID         string  `json:"national_id"`
	ProvinceID         *uint   `json:"province_id"`
	WardUnitID         *uint   `json:"ward_unit_id"`
	AddressLine        string  `json:"address_line"`
	ContactFullName    string  `json:"contact_full_name"`
	ContactPhoneNumber string  `json:"contact_phone_number"`
	ApplicationStatus  string  `json:"application_status"`
	IsPaid             bool    `json:"is_paid"`
	SubmittedAt        *string `json:"submitted_at"`
	CandidateNumber    *string `json:"candidate_number"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type ApplicationAdminFilter struct {
	AdmissionPeriodID *uint  `query:"admission_period_id"`
	ApplicationStatus string `query:"application_status"`
	IsPaid            *bool  `query:"is_paid"`
	ProvinceID        *uint  `query:"province_id"`
	Keyword           string `query:"keyword"`
	Page              int    `query:"page"`
	Limit             int    `query:"limit"`
}

type PaginatedApplicationRes struct {
	Data       []ApplicationRes `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}
