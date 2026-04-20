package dto

type RoleGroupFilter struct {
	Keyword  string `query:"keyword"`
	IsActive *bool  `query:"is_active"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

type RoleGroupCreateReq struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type RoleGroupUpdateReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type RoleGroupStatusReq struct {
	IsActive bool `json:"is_active"`
}

type RoleGroupRes struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type PaginatedRoleGroupRes struct {
	Data       []RoleGroupRes `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type RoleGroupPermissionItem struct {
	PermissionKey   string `json:"permission_key" validate:"required"`
	PermissionValue int64  `json:"permission_value" validate:"min=0"`
}

type RoleGroupPermissionReq struct {
	Permissions []RoleGroupPermissionItem `json:"permissions"`
}

type RoleGroupPermissionRes struct {
	PermissionKey   string `json:"permission_key"`
	PermissionValue int64  `json:"permission_value"`
}
