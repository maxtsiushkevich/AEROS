package rbac

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Action struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

type Resource struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex"`
	Description string
}

type Permission struct {
	ID         uint      `gorm:"primaryKey"`
	ResourceID uint      `gorm:"uniqueIndex:idx_resource_action"`
	ActionID   uint      `gorm:"uniqueIndex:idx_resource_action"`
	Resource   *Resource `gorm:"foreignKey:ResourceID;OnDelete:CASCADE"`
	Action     *Action   `gorm:"foreignKey:ActionID;OnDelete:CASCADE"`
}

type RolePermission struct {
	RoleID       uint
	PermissionID uint
	Role         *Role       `gorm:"foreignKey:RoleID;primaryKey;OnDelete:CASCADE"`
	Permission   *Permission `gorm:"foreignKey:PermissionID;primaryKey;OnDelete:CASCADE"`
}

type UserRole struct {
	UserID uuid.UUID
	RoleID uint
	Role   *Role `gorm:"foreignKey:RoleID;primaryKey;OnDelete:CASCADE"`
}

// DTO
type CreatePermissionRequest struct {
	ResourceID uint `json:"resource_id" binding:"required"`
	ActionID   uint `json:"action_id" binding:"required"`
}

type GrantPermissionRequest struct {
	PermissionID uint `json:"permission_id" binding:"required"`
}
