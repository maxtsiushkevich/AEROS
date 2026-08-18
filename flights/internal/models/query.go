package models

import (
	"time"

	_ "github.com/go-playground/validator/v10"
)

type FlightQuery struct {
	FlightNumber string `validate:"max=8"`
	Origin       string `validate:"max=3"`
	Destination  string `validate:"max=3"`
	Status       string `validate:"omitempty,oneof=Scheduled CheckIn Boarding Delayed Departed Arrived Cancelled Redirected"`
	DateFrom     *time.Time
	DateTo       *time.Time
}
