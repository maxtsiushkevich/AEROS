package storage

import (
	"flights/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Header struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Flight struct {
	Header
	FlightNumber string        `gorm:"size:8;not null"`
	Origin       string        `gorm:"size:3;not null"`
	Destination  string        `gorm:"size:3;not null"`
	Date         time.Time     `gorm:"not null"`
	Status       domain.Status `gorm:"type:status_enum;default:'Scheduled'"`
	Aircraft     string        `gorm:"not null"`
}
