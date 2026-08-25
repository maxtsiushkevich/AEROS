package models

import "github.com/google/uuid"

type UserAuthData struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email          string    `gorm:"unique;not null"`
	HashedPassword string    `gorm:"not null"`
	Version        uint32    `gorm:"default:1"`
}
