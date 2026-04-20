package dto

type AdmissionPeriodFilter struct {
	Keyword    string `query:"keyword"`
	SchoolYear string `query:"school_year"`
	IsOpen     *bool  `query:"is_open"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type AdmissionPeriodReq struct {
	Name       string  `json:"name" validate:"required"`
	SchoolYear string  `json:"school_year" validate:"required"`
	StartDate  *string `json:"start_date"` // Định dạng YYYY-MM-DD
	EndDate    *string `json:"end_date"`   // Định dạng YYYY-MM-DD
	ExamFee    float64 `json:"exam_fee"`
	IsOpen     bool    `json:"is_open"`
	Notes      string  `json:"notes"`
}

type AdmissionPeriodStatusReq struct {
	IsOpen bool `json:"is_open"`
}

type AdmissionPeriodRes struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	SchoolYear string  `json:"school_year"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
	ExamFee    float64 `json:"exam_fee"`
	IsOpen     bool    `json:"is_open"`
	Notes      string  `json:"notes"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type PaginatedAdmissionPeriodRes struct {
	Data       []AdmissionPeriodRes `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}
