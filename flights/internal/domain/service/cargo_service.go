package service

import (
	"context"
	"flights/internal/domain/models"

	"github.com/google/uuid"
)

type CargoService interface {
	CreateCargoManifest(ctx context.Context, flightID uuid.UUID, maxWeightKg int) (*models.CargoManifest, error)
	GetCargoManifest(ctx context.Context, flightID uuid.UUID) (*models.CargoManifest, error)
	UpdateCargoManifest(ctx context.Context, flightID uuid.UUID, maxWeightKg int) (*models.CargoManifest, error)
	DeleteCargoManifest(ctx context.Context, flightID uuid.UUID) error

	CreateCargoItem(ctx context.Context, flightID uuid.UUID, item *models.CargoItem) (*models.CargoItem, error)
	GetCargoItem(ctx context.Context, flightID, itemID uuid.UUID) (*models.CargoItem, error)
	UpdateCargoItem(ctx context.Context, flightID uuid.UUID, item *models.CargoItem) (*models.CargoItem, error)
	DeleteCargoItem(ctx context.Context, flightID, itemID uuid.UUID) error
}
