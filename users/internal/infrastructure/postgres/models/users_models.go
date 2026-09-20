package models

import (
	"time"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Name         string     `gorm:"size:200;not null"`
	Email        string     `gorm:"unique;not null"`
	Birthday     time.Time  `gorm:"not null"`
	MemberNumber string     `gorm:"type:char(10);unique;not null"`
	IsActivated  bool       `gorm:"default:false"`
	ActivatedAt  *time.Time `gorm:"default:null"` // There is no date when the user is not activated
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
