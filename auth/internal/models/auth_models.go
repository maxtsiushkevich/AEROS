package models

import "github.com/google/uuid"

type UserAuth struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email          string    `gorm:"unique;not null"`
	HashedPassword string    `gorm:"not null"`
	Version        uint32    `gorm:"default:1"`
}

func (u *UserAuth) GetVersion() uint32 {
	return u.Version
}

func (UserAuth) TableName() string {
	return "auth"
}

type UserAuthUpdate struct {
	ID                uuid.UUID
	NewEmail          string
	NewHashedPassword string
}
