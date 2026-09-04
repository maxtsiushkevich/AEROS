package models

import "github.com/google/uuid"

type UserAuthData struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email          string    `gorm:"unique;not null"`
	HashedPassword string    `gorm:"not null"`
	Version        uint32    `gorm:"default:1"`
}

func (UserAuthData) TableName() string {
	return "auth"
}

type UserAuthDataUpdate struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	NewEmail    *string   `json:"email" validate:"omitnil,email"`
	NewPassword *string   `json:"password" validate:"omitnil,min=8"`
}
