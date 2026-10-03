package usecase

import (
	"context"
	"flights/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type FlightUseCase interface {
	GetFlights(ctx context.Context, filter *models.FlightFilter) ([]models.Flight, error)
	CreateFlight(ctx context.Context, flightNumber string,
		origin string,
		destination string,
		date time.Time,
		status models.FlightStatus,
		aircraft string) (*models.Flight, error)
	UpdateFlight(ctx context.Context, flight *models.Flight) (*models.Flight, error)
	DeleteFlight(ctx context.Context, id uuid.UUID) error
	CancelFlight(ctx context.Context, id uuid.UUID) error
	RescheduleFlight(ctx context.Context, id uuid.UUID, newDate time.Time) error
	ChangeStatus(ctx context.Context, id uuid.UUID, status string) error
}
