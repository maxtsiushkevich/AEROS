package repository

import (
	"context"

	"flights/internal/domain/models"

	"github.com/google/uuid"
)

type FlightsStorage interface {
	Create(ctx context.Context, flight *models.Flight) (*models.Flight, error)
	Read(ctx context.Context, filter *models.FlightFilter) ([]models.Flight, error)
	ReadById(ctx context.Context, id uuid.UUID) (*models.Flight, error)
	Update(ctx context.Context, flight *models.Flight) (*models.Flight, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
