package usecase

import (
	"context"
	"flights/internal/domain/models"

	"github.com/google/uuid"
)

type CargoUseCase interface {
	GetCargoManifest(ctx context.Context, flightID uuid.UUID) (*models.CargoManifest, error)

	CreateCargoItem(ctx context.Context, flightID uuid.UUID, item *models.CargoItem) (*models.CargoItem, error)
	GetCargoItem(ctx context.Context, flightID, itemID uuid.UUID) (*models.CargoItem, error)
	GetCargoItems(ctx context.Context, flightID uuid.UUID) ([]models.CargoItem, error)
	UpdateCargoItem(ctx context.Context, flightID uuid.UUID, item *models.CargoItem) (*models.CargoItem, error)
	DeleteCargoItem(ctx context.Context, flightID, itemID uuid.UUID) error
}
