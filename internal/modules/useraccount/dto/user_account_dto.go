package dto

type UserAccountFilter struct {
	Keyword  string `query:"keyword"`
	IsActive *bool  `query:"is_active"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

type UserAccountUpdateReq struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	FullName    string  `json:"full_name" validate:"required"`
}

type UserAccountStatusReq struct {
	IsActive bool `json:"is_active"`
}

type UserAccountRes struct {
	ID            uint    `json:"id"`
	Username      string  `json:"username"`
	Email         *string `json:"email"`
	PhoneNumber   *string `json:"phone_number"`
	FullName      string  `json:"full_name"`
	IsActive      bool    `json:"is_active"`
	TotalProfiles int     `json:"total_profiles"` // TODO: Cập nhật count khi có module Hồ Sơ
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type PaginatedUserAccountRes struct {
	Data       []UserAccountRes `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}
