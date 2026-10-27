package postgres

import (
	"time"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Name        string     `gorm:"size:200;not null"`
	Email       string     `gorm:"unique;not null"`
	Birthday    time.Time  `gorm:"not null"`
	IsActivated bool       `gorm:"default:false"`
	ActivatedAt *time.Time `gorm:"default:null"` // There is no date when the user is not activated
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OutboxEvent struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	AggregateType string    `gorm:"type:varchar(255);not null"`
	AggregateID   string    `gorm:"type:varchar(255);not null"`
	EventType     string    `gorm:"type:varchar(255);not null"`
	Payload       []byte    `gorm:"type:jsonb;not null"`
	Attempts      int       `gorm:"not null;default:0"`
	Status        string    `gorm:"type:varchar(50);not null"`
	PartitionKey  *string   `gorm:"type:varchar(255)"`
}
