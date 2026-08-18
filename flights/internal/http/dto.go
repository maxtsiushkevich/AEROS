package http

import (
	"flights/internal/models"
	"time"

	"github.com/google/uuid"
)

type GetFlightsRequestQuery struct {
	FlightNumber string `validate:"max=8"`
	Origin       string `validate:"max=3"`
	Destination  string `validate:"max=3"`
	Status       string `validate:"omitempty,oneof=Scheduled CheckIn Boarding Delayed Departed Arrived Cancelled Redirected"`
	DateFrom     *time.Time
	DateTo       *time.Time
}

type CreateFlightRequest struct {
	FlightNumber string    `json:"flight_number" validate:"required"`
	Origin       string    `json:"origin" validate:"required,len=3"`
	Destination  string    `json:"destination" validate:"required,len=3"`
	Date         time.Time `json:"date" validate:"required,gt=now"`
	Status       string    `json:"status,omitempty" validate:"omitempty,oneof=Scheduled CheckIn Boarding Delayed Departed Arrived Cancelled Redirected"`
	Aircraft     string    `json:"aircraft" validate:"required"`
}

type PatchFlightRequest struct {
	ID           uuid.UUID `json:"id" validate:"required"`
	FlightNumber string    `json:"flight_number"`
	Origin       string    `json:"origin"`
	Destination  string    `json:"destination"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
	Aircraft     string    `json:"aircraft"`
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

type FlightListResponse struct {
	Data []FlightResponse `json:"data"`
}

func (r *GetFlightsRequestQuery) ToServiceQuery() *models.FlightQuery {
	return &models.FlightQuery{
		FlightNumber: r.FlightNumber,
		Origin:       r.Origin,
		Destination:  r.Destination,
		Status:       r.Status,
		DateFrom:     r.DateFrom,
		DateTo:       r.DateTo,
	}
}

func (r *CreateFlightRequest) ToServiceFlight() *models.Flight {
	return &models.Flight{
		FlightNumber: r.FlightNumber,
		Origin:       (r.Origin),
		Destination:  r.Destination,
		Date:         r.Date,
		Status:       models.Status(r.Status),
		Aircraft:     r.Aircraft,
	}
}

func FlightToResponse(flight *models.Flight) FlightResponse {
	return FlightResponse{
		ID:           flight.ID,
		FlightNumber: flight.FlightNumber,
		Origin:       flight.Origin,
		Destination:  flight.Destination,
		Date:         flight.Date,
		Status:       string(flight.Status),
		Aircraft:     flight.Aircraft,
	}
}

func FlightsToResponses(flights []models.Flight) []FlightResponse {
	responses := make([]FlightResponse, len(flights))
	for i, flight := range flights {
		responses[i] = FlightToResponse(&flight)
	}
	return responses
}
