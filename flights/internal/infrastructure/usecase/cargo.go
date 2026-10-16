package usecase

import (
	"context"
	appusecase "flights/internal/application/usecase"
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

func (c *CargoUseCase) GetCargoManifest(ctx context.Context, cargoId uuid.UUID) (*models.CargoManifest, error) {
	manifest, err := c.cargoStorage.Read(ctx, cargoId)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (c *CargoUseCase) AddCargoItem(ctx context.Context, req appusecase.AddCargoItemRequest) (*models.CargoItem, error) {
	manifest, err := c.cargoStorage.Read(ctx, req.ManifestID)
	if err != nil {
		return nil, err
	}

	item, err := models.NewCargoItem(req.ManifestID, models.CargoType(req.CargoType), req.WeightKg, req.PassengerID, req.Description)
	if err != nil {
		return nil, err
	}

	err = manifest.AddCargoItem(*item)
	if err != nil {
		return nil, err
	}

	err = c.cargoStorage.Update(ctx, manifest)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (c *CargoUseCase) GetCargoItem(ctx context.Context, req appusecase.GetCargoItemRequest) (*models.CargoItem, error) {
	manifest, err := c.cargoStorage.Read(ctx, req.ManifestID)
	if err != nil {
		return nil, err
	}

	item, err := manifest.GetItem(req.ItemID)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *CargoUseCase) GetCargoItems(ctx context.Context, cargoId uuid.UUID) ([]models.CargoItem, error) {
	manifest, err := c.cargoStorage.Read(ctx, cargoId)
	if err != nil {
		return nil, err
	}

	return manifest.Items(), nil
}

func (c *CargoUseCase) DeleteCargoItem(ctx context.Context, cargoId, cargoItemId uuid.UUID) error {
	manifest, err := c.cargoStorage.Read(ctx, cargoId)
	if err != nil {
		return err
	}

	manifest.DeleteCargoItem(cargoItemId)

	err = c.cargoStorage.Update(ctx, manifest)
	if err != nil {
		return err
	}

	err = c.cargoStorage.DeleteCargoItem(ctx, cargoItemId)
	if err != nil {
		return err
	}

	return nil
}

func (c *CargoUseCase) MoveCargoItem(ctx context.Context, req appusecase.MoveCargoItemRequest) error {
	if req.CargoID == req.NewCargoID {
		return nil
	}

	manifest, err := c.cargoStorage.Read(ctx, req.CargoID)
	if err != nil {
		return err
	}

	item, err := manifest.GetItem(req.ItemID)
	if err != nil {
		return err
	}

	newManifest, err := c.cargoStorage.Read(ctx, req.NewCargoID)
	if err != nil {
		return err
	}

	if err = c.service.MoveItem(ctx, &item, manifest, newManifest); err != nil {
		return err
	}

	if err = c.cargoStorage.Update(ctx, manifest); err != nil {
		return err
	}

	if err = c.cargoStorage.Update(ctx, newManifest); err != nil {
		return err
	}

	return nil
}
