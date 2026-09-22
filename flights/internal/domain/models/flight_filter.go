package models

import (
	"time"

	_ "github.com/go-playground/validator/v10"
)

type FlightFilter struct {
	FlightNumber *string `validate:"omitnil,max=8"`
	Origin       *string `validate:"max=3"`
	Destination  *string `validate:"max=3"`
	Status       *string `validate:"omitempty,oneof=Scheduled CheckIn Boarding Delayed Departed Arrived Cancelled Redirected"`
	DateFrom     *time.Time
	DateTo       *time.Time
	Limit        int
	Page         int
}

// type FlightUpdate struct {
// 	ID           uuid.UUID
// 	FlightNumber *string
// 	Origin       *string
// 	Destination  *string
// 	Date         *time.Time
// 	Status       *models.FlightStatus
// 	Aircraft     *string
// }
