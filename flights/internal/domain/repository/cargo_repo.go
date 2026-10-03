package repository

import (
	"context"

	"flights/internal/domain/models"

	"github.com/google/uuid"
)

type CargoRepository interface {
	Create(ctx context.Context, manifest *models.CargoManifest) error
	FindByFlightID(ctx context.Context, flightID uuid.UUID) (*models.CargoManifest, error)
	Update(ctx context.Context, manifest *models.CargoManifest) error
	Delete(ctx context.Context, flightID uuid.UUID) error
}
