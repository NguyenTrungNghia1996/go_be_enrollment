package entity

import "time"

type RoleGroup struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"uniqueIndex;size:50;not null"`
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AdminUserRoleGroup struct {
	AdminUserID uint `gorm:"primaryKey"`
	RoleGroupID uint `gorm:"primaryKey"`
	CreatedAt   time.Time
}

type RoleGroupPermission struct {
	RoleGroupID     uint   `gorm:"primaryKey"`
	PermissionKey   string `gorm:"primaryKey;size:100"`
	PermissionValue int64  `gorm:"not null;default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Menu struct {
	ID            uint   `gorm:"primaryKey"`
	ParentID      uint   `gorm:"index;default:0"`
	Title         string `gorm:"size:100;not null"`
	MenuKey       string `gorm:"size:100;uniqueIndex"`
	Icon          string `gorm:"size:100"`
	Url           string `gorm:"size:255;not null"`
	PermissionBit int    `gorm:"default:0"`
	IsActive      bool   `gorm:"default:true"`
	SortOrder     int    `gorm:"default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableNames matching the initialized schema
func (RoleGroup) TableName() string { return "role_groups" }
func (AdminUserRoleGroup) TableName() string { return "admin_user_role_groups" }
func (RoleGroupPermission) TableName() string { return "role_group_permissions" }
func (Menu) TableName() string { return "menus" }
