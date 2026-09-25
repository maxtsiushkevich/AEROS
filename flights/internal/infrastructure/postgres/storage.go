package postgres

import (
	"context"
	"errors"
	"flights/internal/config"
	domain_err "flights/internal/domain/errors"
	"flights/internal/domain/models"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresFlightsStorage struct {
	config *config.Config
	logger *slog.Logger
	db     *gorm.DB
}

func CreateStorage(cfg *config.Config, l *slog.Logger) *PostgresFlightsStorage {
	return &PostgresFlightsStorage{
		config: cfg,
		logger: l,
	}
}

func (s *PostgresFlightsStorage) Open() error {
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

	pool, _ := s.db.DB()

	pool.SetMaxOpenConns(5)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(30 * time.Second)
	pool.SetConnMaxIdleTime(15 * time.Second)

	return nil
}

func (s *PostgresFlightsStorage) autoMigrateModels() error {
	if err := s.db.AutoMigrate(&Flight{}); err != nil {
		s.logger.Error("Failed to auto-migrate models", "err", err)
		return err
	}

	return nil
}

func (s *PostgresFlightsStorage) Close() error {
	database, err := s.db.DB()
	if err != nil {
		return err
	}
	database.Close()
	s.logger.Info("Database connection closed", "err", err)
	return nil
}

func (s *PostgresFlightsStorage) Create(ctx context.Context, flight *models.Flight) (*models.Flight, error) {
	err := s.db.WithContext(ctx).Create(flight).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create flight in db: %w", err)
	}
	return flight, nil
}

func (s *PostgresFlightsStorage) Read(ctx context.Context, filter *models.FlightFilter) ([]models.Flight, error) {
	var flights []models.Flight
	query := s.db.WithContext(ctx)

	if filter == nil {
		filter = &models.FlightFilter{}
	}

	if filter.FlightNumber != nil {
		query = query.Where("flight_number = ?", filter.FlightNumber)
	}

	if filter.Origin != nil {
		query = query.Where("origin = ?", filter.Origin)
	}

	if filter.Destination != nil {
		query = query.Where("destination = ?", filter.Destination)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.DateFrom != nil {
		query = query.Where("date >= ?", filter.DateFrom)
	}

	if filter.DateTo != nil {
		query = query.Where("date <= ?", filter.DateTo)
	}

	if filter.Limit > 0 {
		query = query.Order("date ASC, id ASC").Limit(filter.Limit)
		if filter.Page > 1 {
			query = query.Offset((filter.Page - 1) * filter.Limit)
		}
	}

	if err := query.Find(&flights).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain_err.ErrFlightNotFound
		}
		s.logger.Error("Failed to read flights", "err", err)
		return nil, err
	}

	return flights, nil
}

func (s *PostgresFlightsStorage) Update(ctx context.Context, flight *models.Flight) (*models.Flight, error) {
	if flight == nil {
		return nil, fmt.Errorf("flight is nil")
	}

	var existing Flight
	if err := s.db.WithContext(ctx).First(&existing, "id = ?", flight.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_err.ErrFlightNotFound
		}
		return nil, fmt.Errorf("failed to fetch flight for update: %w", err)
	}

	DomainToDB(flight, &existing)

	result := s.db.WithContext(ctx).
		Model(&Flight{}).
		Where("id = ?", flight.Id).
		Updates(existing)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update flight in db: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, domain_err.ErrFlightNotFound
	}

	return ToDomain(&existing), nil
}

func (s *PostgresFlightsStorage) Delete(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("id = ?", id).Delete(&Flight{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete flight in db: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain_err.ErrFlightNotFound
	}

	return nil
}

func (s *PostgresFlightsStorage) ReadById(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	var dbFlight Flight
	if err := s.db.WithContext(ctx).First(&dbFlight, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain_err.ErrFlightNotFound
		}
		s.logger.Error("Flight not found", "id", id)
	}
	return ToDomain(&dbFlight), nil
}
