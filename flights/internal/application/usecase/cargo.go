package usecase

import (
	"context"
	"flights/internal/domain/models"

	"github.com/google/uuid"
)

type GetCargoManifestRequest struct {
	ID uuid.UUID
}

type AddCargoItemRequest struct {
	ManifestID  uuid.UUID
	CargoType   string
	WeightKg    int
	PassengerID *uuid.UUID
	Description string
}

type GetCargoItemRequest struct {
	FlightID uuid.UUID
	ItemID   uuid.UUID
}

type GetCargoItemsRequest struct {
	FlightID uuid.UUID
}

type DeleteCargoItemRequest struct {
	FlightID uuid.UUID
	ItemID   uuid.UUID
}

type MoveCargoItemRequest struct {
	ItemID      uuid.UUID
	NewFlightID uuid.UUID
}

type CargoUseCase interface {
	GetCargoManifest(ctx context.Context, req GetCargoManifestRequest) (*models.CargoManifest, error)
	AddCargoItem(ctx context.Context, req AddCargoItemRequest) (*models.CargoItem, error)
	GetCargoItem(ctx context.Context, req GetCargoItemRequest) (*models.CargoItem, error)
	GetCargoItems(ctx context.Context, req GetCargoItemsRequest) ([]models.CargoItem, error)
	DeleteCargoItem(ctx context.Context, req DeleteCargoItemRequest) error
	MoveCargoItem(ctx context.Context, req MoveCargoItemRequest) error
}
