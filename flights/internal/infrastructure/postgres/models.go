package postgres

import (
	"flights/internal/domain/models"
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
	FlightNumber string              `gorm:"size:8;not null"`
	Origin       string              `gorm:"size:3;not null"`
	Destination  string              `gorm:"size:3;not null"`
	Date         time.Time           `gorm:"not null"`
	Status       models.FlightStatus `gorm:"type:status_enum;default:'Scheduled'"`
	Aircraft     string              `gorm:"not null"`
}

type CargoManifest struct {
	Header
	FlightID        uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex"`
	MaxWeightKg     int         `gorm:"not null"`
	CurrentWeightKg int         `gorm:"not null;default:0"`
	CargoItems      []CargoItem `gorm:"foreignKey:ManifestID;references:ID"`

	Flight Flight `gorm:"foreignKey:FlightID;references:ID"`
}

type CargoItem struct {
	ID          uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ManifestID  uuid.UUID        `gorm:"type:uuid;not null;index"`
	CargoType   models.CargoType `gorm:"type:cargo_type_enum;not null"`
	WeightKg    int              `gorm:"not null;default:0"`
	PassengerID *uuid.UUID       `gorm:"type:uuid"`
	Description string           `gorm:"not null"`
	PackedAt    time.Time        `gorm:"not null"`

	Manifest CargoManifest `gorm:"foreignKey:ManifestID;references:ID"`
}
