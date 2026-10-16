package service

import (
	"context"
	"flights/internal/domain/models"
	"flights/internal/domain/repository"
)

type CargoDomainService struct {
	repo repository.CargoRepository
}

func NewCargoDomainService(repo repository.CargoRepository) *CargoDomainService {
	return &CargoDomainService{repo: repo}
}

func (s CargoDomainService) MoveItem(ctx context.Context, item *models.CargoItem, manifest, newManifest *models.CargoManifest) error {
	if manifest == nil || newManifest == nil {
		return nil
	}
	if manifest.ID == newManifest.ID {
		return nil
	}

	manifest.DeleteCargoItem(item.ID)

	item.Move(newManifest.ID)

	err := newManifest.AddCargoItem(*item)
	if err != nil {
		// return to original manifest
		item.Move(manifest.ID)
		return err
	}

	return nil
}
