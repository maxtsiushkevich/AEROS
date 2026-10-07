package postgres

import (
	"context"
	"errors"
	domain_err "flights/internal/domain/errors"
	"flights/internal/domain/models"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresFlightsStorage struct {
	logger *slog.Logger
	db     *gorm.DB
}

func CreateFlightsStorage(db *gorm.DB, l *slog.Logger) *PostgresFlightsStorage {
	return &PostgresFlightsStorage{
		db:     db,
		logger: l,
	}
}

func (s *PostgresFlightsStorage) Create(ctx context.Context, flight *models.Flight) (*models.Flight, error) {
	if flight == nil {
		return nil, domain_err.ErrFlightNotFound
	}

	dbFlight := FlightFromDomain(flight)
	if err := s.db.WithContext(ctx).Create(dbFlight).Error; err != nil {
		return nil, fmt.Errorf("failed to create flight in db: %w", err)
	}

	return FlightToDomain(dbFlight), nil
}

func (s *PostgresFlightsStorage) Read(ctx context.Context, filter *models.FlightFilter) ([]models.Flight, error) {
	var dbFlights []Flight
	query := s.db.WithContext(ctx).Model(&Flight{})

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

	if err := query.Find(&dbFlights).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_err.ErrFlightNotFound
		}
		s.logger.Error("Failed to read flights", "err", err)
		return nil, err
	}

	return FlightsToDomain(dbFlights), nil
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

	FlightToDb(flight, &existing)

	result := s.db.WithContext(ctx).
		Model(&Flight{}).
		Where("id = ?", flight.Id).
		Updates(&existing)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update flight in db: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, domain_err.ErrFlightNotFound
	}

	return FlightToDomain(&existing), nil
}

func (s *PostgresFlightsStorage) Delete(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("id = ?", id).Delete(&Flight{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete flight in db: %w", result.Error)
	}

	return nil
}

func (s *PostgresFlightsStorage) ReadById(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	var dbFlight Flight
	if err := s.db.WithContext(ctx).First(&dbFlight, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_err.ErrFlightNotFound
		}
		s.logger.Error("Flight not found", "id", id, "err", err)
		return nil, err
	}
	return FlightToDomain(&dbFlight), nil
}
