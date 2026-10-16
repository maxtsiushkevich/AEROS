package service

import (
	"context"
	"flights/internal/domain/models"
)

type CargoService interface {
	MoveItem(ctx context.Context, item *models.CargoItem, manifest, newManifest *models.CargoManifest) error
}
