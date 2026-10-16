package usecase

import (
	"context"
	"flights/internal/domain/models"

	"github.com/google/uuid"
)

type AddCargoItemRequest struct {
	ManifestID  uuid.UUID
	CargoType   string
	WeightKg    int
	PassengerID *uuid.UUID
	Description string
}

type GetCargoItemRequest struct {
	ManifestID uuid.UUID
	ItemID     uuid.UUID
}

type MoveCargoItemRequest struct {
	ItemID     uuid.UUID
	CargoID    uuid.UUID
	NewCargoID uuid.UUID
}

type CargoUseCase interface {
	GetCargoManifest(ctx context.Context, cargoId uuid.UUID) (*models.CargoManifest, error)
	AddCargoItem(ctx context.Context, req AddCargoItemRequest) (*models.CargoItem, error)
	GetCargoItem(ctx context.Context, req GetCargoItemRequest) (*models.CargoItem, error)
	GetCargoItems(ctx context.Context, cargoId uuid.UUID) ([]models.CargoItem, error)
	DeleteCargoItem(ctx context.Context, cargoId, cargoItemId uuid.UUID) error
	MoveCargoItem(ctx context.Context, req MoveCargoItemRequest) error
}
