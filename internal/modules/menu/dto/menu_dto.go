package dto

type MenuCreateReq struct {
	ParentID      uint   `json:"parent_id"`
	Title         string `json:"title" validate:"required"`
	MenuKey       string `json:"menu_key" validate:"required"`
	Icon          string `json:"icon"`
	Url           string `json:"url" validate:"required"`
	PermissionBit int    `json:"permission_bit"`
	IsActive      bool   `json:"is_active"`
	SortOrder     int    `json:"sort_order"`
}

type MenuUpdateReq struct {
	ParentID      uint   `json:"parent_id"`
	Title         string `json:"title" validate:"required"`
	MenuKey       string `json:"menu_key" validate:"required"`
	Icon          string `json:"icon"`
	Url           string `json:"url" validate:"required"`
	PermissionBit int    `json:"permission_bit"`
	IsActive      bool   `json:"is_active"`
	SortOrder     int    `json:"sort_order"`
}

type MenuRes struct {
	ID            uint      `json:"id"`
	ParentID      uint      `json:"parent_id"`
	Title         string    `json:"title"`
	MenuKey       string    `json:"menu_key"`
	Icon          string    `json:"icon"`
	Url           string    `json:"url"`
	PermissionBit int       `json:"permission_bit"`
	IsActive      bool      `json:"is_active"`
	SortOrder     int       `json:"sort_order"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	Children      []MenuRes `json:"children"`
}
