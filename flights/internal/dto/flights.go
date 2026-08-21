package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetFlightsRequestQuery struct {
	FlightNumber string `validate:"max=8"`
	Origin       string `validate:"omitempty,max=3,alpha"`
	Destination  string `validate:"omitempty,max=3,alpha"`
	Status       string `validate:"omitempty,oneof=Scheduled CheckIn Boarding Delayed Departed Arrived Cancelled Redirected"`
	DateFrom     *time.Time
	DateTo       *time.Time
}

type CreateFlightRequest struct {
	FlightNumber *string    `json:"flight_number" validate:"required,max=8"`
	Origin       *string    `json:"origin" validate:"required,len=3,alpha"`
	Destination  *string    `json:"destination" validate:"required,len=3,alpha"`
	Date         *time.Time `json:"date" validate:"required,gt=now"`
	Status       *string    `json:"status,omitempty" validate:"omitempty,oneof=Scheduled CheckIn Boarding Delayed Departed Arrived Cancelled Redirected"`
	Aircraft     *string    `json:"aircraft" validate:"required"`
}

type PatchFlightRequest struct {
	ID           uuid.UUID  `json:"id" validate:"required"`
	FlightNumber *string    `json:"flight_number,omitempty" validate:"omitempty,max=8"`
	Origin       *string    `json:"origin,omitempty" validate:"omitempty,len=3,alpha"`
	Destination  *string    `json:"destination,omitempty" validate:"omitempty,len=3,alpha"`
	Date         *time.Time `json:"date,omitempty" validate:"omitempty"`
	Status       *string    `json:"status,omitempty" validate:"omitempty"`
	Aircraft     *string    `json:"aircraft,omitempty" validate:"omitempty"`
}

type FlightResponse struct {
	ID           uuid.UUID `json:"id"`
	FlightNumber string    `json:"flight_number"`
	Origin       string    `json:"origin"`
	Destination  string    `json:"destination"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
	Aircraft     string    `json:"aircraft"`
}
