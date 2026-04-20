package dto

type ProvinceFilter struct {
	Keyword  string `query:"keyword"`
	IsActive *bool  `query:"is_active"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

type ProvinceCreateReq struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type ProvinceUpdateReq struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type ProvinceStatusReq struct {
	IsActive bool `json:"is_active"`
}

type ProvinceRes struct {
	ID        uint   `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PaginatedProvinceRes struct {
	Data       []ProvinceRes `json:"data"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}
