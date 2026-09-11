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
	Name      string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Resource struct {
	Name        string `gorm:"primaryKey"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ResourceName string    `gorm:"primaryKey;column:resource_name"`
	ActionName   string    `gorm:"primaryKey;column:action_name"`
	Resource     *Resource `gorm:"foreignKey:ResourceName;constraint:OnDelete:CASCADE;"`
	Action       *Action   `gorm:"foreignKey:ActionName;constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RolePermission struct {
	RoleName     string      `gorm:"primaryKey"`
	ResourceName string      `gorm:"primaryKey"`
	ActionName   string      `gorm:"primaryKey"`
	Role         *Role       `gorm:"foreignKey:RoleName;constraint:OnDelete:CASCADE;"`
	Permission   *Permission `gorm:"foreignKey:ResourceName,ActionName;constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRole struct {
	UserID    uuid.UUID `gorm:"primaryKey"`
	RoleName  string    `gorm:"primaryKey"`
	Role      *Role     `gorm:"foreignKey:RoleName;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
