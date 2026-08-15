package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string

const (
	Scheduled  Status = "Scheduled"
	CheckIn    Status = "CheckIn"
	Boarding   Status = "Boarding"
	Delayed    Status = "Delayed"
	Departed   Status = "Departed"
	Arrived    Status = "Arrived"
	Cancelled  Status = "Cancelled"
	Redirected Status = "Redirected"
)

type Model struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Flight struct {
	Model
	FlightNumber string    `gorm:"size:8;not null"`
	Origin       string    `gorm:"size:3;not null"`
	Destination  string    `gorm:"size:3;not null"`
	Date         time.Time `gorm:"not null"`
	Status       Status    `gorm:"type:status_enum;default:'Scheduled'"`
	Aircraft     string    `gorm:"not null"`
}
