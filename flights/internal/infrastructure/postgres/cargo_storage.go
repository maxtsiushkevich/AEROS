package postgres

import (
	"context"
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
	return nil
}

func (s *PostgresCargoStorage) FindByFlightID(ctx context.Context, flightID uuid.UUID) (*models.CargoManifest, error) {
	return nil, nil
}

func (s *PostgresCargoStorage) Update(ctx context.Context, manifest *models.CargoManifest) error {
	return nil
}

func (s *PostgresCargoStorage) Delete(ctx context.Context, flightID uuid.UUID) error {
	return nil
}
