package rbac

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Name        string `gorm:"primaryKey"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Action struct {
	Name string `gorm:"primaryKey"`
}

type Resource struct {
	Name        string `gorm:"primaryKey"`
	Description string
}

type Permission struct {
	ResourceName string    `gorm:"primaryKey;column:resource_name"`
	ActionName   string    `gorm:"primaryKey;column:action_name"`
	Resource     *Resource `gorm:"foreignKey:ResourceName;OnDelete:CASCADE"`
	Action       *Action   `gorm:"foreignKey:ActionName;OnDelete:CASCADE"`
}

type RolePermission struct {
	RoleName     string      `gorm:"primaryKey"`
	ResourceName string      `gorm:"primaryKey"`
	ActionName   string      `gorm:"primaryKey"`
	Role         *Role       `gorm:"foreignKey:RoleName;OnDelete:CASCADE"`
	Permission   *Permission `gorm:"foreignKey:ResourceName,ActionName;OnDelete:CASCADE"`
}

type UserRole struct {
	UserID   uuid.UUID `gorm:"primaryKey"`
	RoleName string    `gorm:"primaryKey"`
	Role     *Role     `gorm:"foreignKey:RoleName;OnDelete:CASCADE"`
}

// DTOs
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateActionRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateResourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreatePermissionRequest struct {
	ResourceName string `json:"resource_name" binding:"required"`
	ActionName   string `json:"action_name" binding:"required"`
}

type GrantPermissionRequest struct {
	RoleName     string `json:"role_name" binding:"required"`
	ResourceName string `json:"resource_name" binding:"required"`
	ActionName   string `json:"action_name" binding:"required"`
}

type RevokePermissionRequest struct {
	RoleName     string `json:"role_name" binding:"required"`
	ResourceName string `json:"resource_name" binding:"required"`
	ActionName   string `json:"action_name" binding:"required"`
}

type AssignRoleRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	RoleName string    `json:"role_name" binding:"required"`
}

type RemoveRoleRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	RoleName string    `json:"role_name" binding:"required"`
}
