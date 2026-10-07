package usecase

import (
	"context"
	"flights/internal/domain/models"
	"flights/internal/domain/repository"
	"flights/internal/domain/service"
	domService "flights/internal/infrastructure/service"

	"github.com/google/uuid"
)

type CargoUseCase struct {
	cargoStorage repository.CargoRepository
	service      service.CargoService
}

func NewCargoUseCase(storage repository.CargoRepository) *CargoUseCase {
	return &CargoUseCase{
		cargoStorage: storage,
		service:      domService.NewCargoDomainService(storage),
	}
}

func (c *CargoUseCase) GetCargoManifest(ctx context.Context, flightID uuid.UUID) (*models.CargoManifest, error) {
	manifest, err := c.cargoStorage.FindByFlightID(ctx, flightID)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (c *CargoUseCase) CreateCargoItem(ctx context.Context, flightID uuid.UUID, item *models.CargoItem) (*models.CargoItem, error) {
	return nil, nil
}

func (c *CargoUseCase) GetCargoItem(ctx context.Context, flightID, itemID uuid.UUID) (*models.CargoItem, error) {
	return nil, nil
}

func (c *CargoUseCase) GetCargoItems(ctx context.Context, flightID uuid.UUID) ([]models.CargoItem, error) {
	return nil, nil
}

func (c *CargoUseCase) DeleteCargoItem(ctx context.Context, flightID, itemID uuid.UUID) error {
	return nil
}

func (c *CargoUseCase) MoveCargoItem(ctx context.Context, itemID, newFlightId uuid.UUID) error {
	return nil
}
