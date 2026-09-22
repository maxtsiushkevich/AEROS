package postgres

import (
	"context"
	"flights/internal/config"
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
	// return nil, nil
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
		s.logger.Error("Failed to read flights", "err", err)
		return nil, err
	}

	return flights, nil
}

func (s *PostgresFlightsStorage) Update(ctx context.Context, flight *models.Flight) (*models.Flight, error) {
	return nil, nil

	// if flight == nil {
	// 	return nil, errors_pkg.New("flight update data is nil")
	// }

	// var existing Flight
	// if err := s.db.WithContext(ctx).First(&existing, "id = ?", flight.ID).Error; err != nil {
	// 	if errors_pkg.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, errors.FlightNotFoundError
	// 	}
	// 	return nil, fmt.Errorf("failed to fetch flight for update: %w", err)
	// }

	// if flight.FlightNumber != nil {
	// 	existing.FlightNumber = *flight.FlightNumber
	// }
	// if flight.Origin != nil {
	// 	existing.Origin = *flight.Origin
	// }
	// if flight.Destination != nil {
	// 	existing.Destination = *flight.Destination
	// }
	// if flight.Date != nil {
	// 	existing.Date = *flight.Date
	// }
	// if flight.Status != nil {
	// 	existing.Status = *flight.Status
	// }
	// if flight.Aircraft != nil {
	// 	existing.Aircraft = *flight.Aircraft
	// }

	// if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
	// 	return nil, fmt.Errorf("failed to update flight in db: %w", err)
	// }

	// return &existing, nil
}

func (s *PostgresFlightsStorage) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
	// result := s.db.WithContext(ctx).Where("id = ?", id).Delete(&Flight{})
	// if result.Error != nil {
	// 	return fmt.Errorf("failed to delete flight in db: %w", result.Error)
	// }
	// if result.RowsAffected == 0 {
	// 	return errors.FlightNotFoundError
	// }

	// return nil
}

func (s *PostgresFlightsStorage) ReadById(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	return nil, nil
}
