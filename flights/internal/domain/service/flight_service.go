package service

import (
	"context"
	"flights/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type FlightService interface {
	GetFlights(ctx context.Context, filter *models.FlightFilter) ([]models.Flight, error)
	CreateFlight(ctx context.Context, flightNumber string,
		origin string,
		destination string,
		date time.Time,
		status models.FlightStatus,
		aircraft string) (*models.Flight, error)
	UpdateFlight(ctx context.Context, flight *models.Flight) (*models.Flight, error)
	DeleteFlight(ctx context.Context, id uuid.UUID) error
}
