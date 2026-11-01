package outbox

import (
	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	AggregateType string    `gorm:"type:varchar(255);not null"`
	AggregateID   string    `gorm:"type:varchar(255);not null"`
	EventType     string    `gorm:"type:varchar(255);not null"`
	Payload       []byte    `gorm:"type:jsonb;not null"`
	Attempts      int       `gorm:"not null;default:0"`
	Status        string    `gorm:"type:varchar(50);not null;index"`
	PartitionKey  *string   `gorm:"type:varchar(255)"`
}

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusDelivered  = "delivered"
)

func (OutboxEvent) TableName() string {
	return "outbox_events"
}
