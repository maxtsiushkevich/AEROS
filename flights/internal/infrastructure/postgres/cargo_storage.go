package postgres

import (
	"context"
	"errors"
	domain_err "flights/internal/domain/errors"
	"flights/internal/domain/models"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresCargoStorage struct {
	logger *slog.Logger
	db     *gorm.DB
}

func CreateCargoStorage(db *gorm.DB, l *slog.Logger) *PostgresCargoStorage {
	return &PostgresCargoStorage{
		db:     db,
		logger: l,
	}
}

func (s *PostgresCargoStorage) Create(ctx context.Context, manifest *models.CargoManifest) error {
	if manifest == nil {
		return domain_err.ErrCargoManifestNotFound
	}

	dbManifest := CargoManifestFromDomain(manifest)
	err := gorm.G[CargoManifest](s.db).Create(ctx, &dbManifest)
	return err
}

func (s *PostgresCargoStorage) Read(ctx context.Context, id uuid.UUID) (*models.CargoManifest, error) {
	cm, err := gorm.G[CargoManifest](s.db).Preload("CargoItems", nil).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_err.ErrCargoManifestNotFound
		}
		return nil, err
	}

	return cm.ToDomain(), nil
}

func (s *PostgresCargoStorage) FindByFlightID(ctx context.Context, flightID uuid.UUID) (*models.CargoManifest, error) {
	manifest, err := gorm.G[CargoManifest](s.db).Preload("CargoItems", nil).Where("flight_id = ?", flightID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_err.ErrCargoManifestNotFound
		}
		return nil, err
	}

	return manifest.ToDomain(), nil
}

func (s *PostgresCargoStorage) Update(ctx context.Context, manifest *models.CargoManifest) error {
	if manifest == nil {
		return domain_err.ErrCargoManifestNotFound
	}
	dbManifest := CargoManifestFromDomain(manifest)
	_, err := gorm.G[CargoManifest](s.db).Where("id = ?", manifest.ID).Updates(ctx, dbManifest)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresCargoStorage) Delete(ctx context.Context, flightID uuid.UUID) error {
	_, err := gorm.G[CargoManifest](s.db).Where("flight_id = ?", flightID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
