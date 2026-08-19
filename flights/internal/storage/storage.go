package storage

import (
	"context"
	"flights/internal/config"
	"flights/internal/models"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FlightsStorage interface {
	Open() error
	Close() error
	Create(ctx context.Context, flight *models.Flight) error
	Read(ctx context.Context, filter *models.FlightQuery) ([]models.Flight, error)
	Update(ctx context.Context, flight *models.Flight) (*models.Flight, error)
	Delete(ctx context.Context, id uuid.UUID)
}

type FlightsPostgresStorage struct {
	config *config.Config
	logger *slog.Logger
	db     *gorm.DB
}

func CreateStorage(cfg *config.Config, l *slog.Logger) FlightsStorage {
	return &FlightsPostgresStorage{
		config: cfg,
		logger: l,
	}
}

func (s *FlightsPostgresStorage) Open() error {
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s",
		s.config.Database.User,
		s.config.Database.Password,
		s.config.Database.Address,
		s.config.Database.DbName)

	var err error

	s.db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
	s.logger.Debug("Open db connect", "connString", connString)
	if err != nil {
		return err
	}

	if err := s.autoMigrateModels(); err != nil {
		return err
	}
	s.logger.Debug("Init schemas")

	return nil
}

func (s *FlightsPostgresStorage) autoMigrateModels() error {
	if err := s.db.AutoMigrate(&models.Flight{}); err != nil {
		s.logger.Error("Failed to auto-migrate models", "err", err)
		return err
	}

	return nil
}

func (s *FlightsPostgresStorage) Close() error {
	database, err := s.db.DB()
	if err != nil {
		return err
	}
	database.Close()
	s.logger.Info("Database connection closed", "err", err)
	return nil
}

func (s *FlightsPostgresStorage) Create(ctx context.Context, flight *models.Flight) error {
	err := s.db.WithContext(ctx).Create(flight).Error
	if err != nil {
		return fmt.Errorf("failed to create flight in db: %w", err)
	}
	return nil
}

func (s *FlightsPostgresStorage) Read(ctx context.Context, filter *models.FlightQuery) ([]models.Flight, error) {
	var flights []models.Flight
	query := s.db.WithContext(ctx)

	if filter == nil {
		filter = &models.FlightQuery{}
	}

	if filter.FlightNumber != "" {
		query = query.Where("flight_number = ?", filter.FlightNumber)
	}

	if filter.Origin != "" {
		query = query.Where("origin = ?", filter.Origin)
	}

	if filter.Destination != "" {
		query = query.Where("destination = ?", filter.Destination)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.DateFrom != nil {
		query = query.Where("date >= ?", filter.DateFrom)
	}

	if filter.DateTo != nil {
		query = query.Where("date <= ?", filter.DateTo)
	}

	if err := query.Find(&flights).Error; err != nil {
		s.logger.Error("Failed to read flights", "err", err)
		return nil, err
	}

	return flights, nil
}

func (s *FlightsPostgresStorage) Update(ctx context.Context, flight *models.Flight) (*models.Flight, error) {
	return &models.Flight{}, nil
}

func (s *FlightsPostgresStorage) Delete(ctx context.Context, id uuid.UUID) {

}
