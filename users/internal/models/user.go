package models

import (
	"time"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Name         string
	Email        string
	Age          string
	MemberNumber uint64
	IsActivated  bool
	ActivatedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
