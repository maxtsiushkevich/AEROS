package usecase

import (
	"context"
	appusecase "flights/internal/application/usecase"
	"flights/internal/domain/models"
	"flights/internal/domain/repository"
	"flights/internal/domain/service"
	domService "flights/internal/infrastructure/service"
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

func (c *CargoUseCase) GetCargoManifest(ctx context.Context, req appusecase.GetCargoManifestRequest) (*models.CargoManifest, error) {
	manifest, err := c.cargoStorage.Read(ctx, req.ID)
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
	return nil, nil
}

func (c *CargoUseCase) GetCargoItems(ctx context.Context, req appusecase.GetCargoItemsRequest) ([]models.CargoItem, error) {
	return nil, nil
}

func (c *CargoUseCase) DeleteCargoItem(ctx context.Context, req appusecase.DeleteCargoItemRequest) error {
	return nil
}

func (c *CargoUseCase) MoveCargoItem(ctx context.Context, req appusecase.MoveCargoItemRequest) error {
	return nil
}
