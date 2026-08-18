package http

import (
	"time"

	"github.com/google/uuid"
)

type FlightResponse struct {
	ID           uuid.UUID `json:"id"`
	FlightNumber string    `json:"flight_number"`
	Origin       string    `json:"origin"`
	Destination  string    `json:"destination"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
	Aircraft     string    `json:"aircraft"`
}

type FlightListResponse struct {
	Data []FlightResponse `json:"data"`
}

type CreateFlightRequest struct {
	FlightNumber string    `json:"flight_number" validate:"required"`
	Origin       string    `json:"origin" validate:"required,len=3"`
	Destination  string    `json:"destination" validate:"required,len=3"`
	Date         time.Time `json:"date" validate:"required,gt=now"`
	Status       string    `json:"status" validate:"required,oneof=Scheduled CheckIn Boarding Delayed Departed Arrived Cancelled Redirected"`
	Aircraft     string    `json:"aircraft" validate:"required"`
}

type FlightQueryParams struct {
	FlightNumber string `validate:"max=8"`
	Origin       string `validate:"max=3"`
	Destination  string `validate:"max=3"`
	Status       string `validate:"omitempty,oneof=Scheduled CheckIn Boarding Delayed Departed Arrived Cancelled Redirected"`
	DateFrom     *time.Time
	DateTo       *time.Time
}
