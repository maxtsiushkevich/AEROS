package models

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	HashedPassword string
	Version        uint32 `gorm:"default:1"`
}
