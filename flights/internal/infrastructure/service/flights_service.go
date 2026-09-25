package service

import (
	"context"
	"flights/internal/domain/models"
	"flights/internal/domain/repository"
	"time"

	"github.com/google/uuid"
)

type FlightService struct {
	storage repository.FlightsStorage
}

func NewFlightService(storage repository.FlightsStorage) *FlightService {
	return &FlightService{
		storage: storage,
	}
}

func (s *FlightService) GetFlights(ctx context.Context, filter *models.FlightFilter) ([]models.Flight, error) {
	flights, err := s.storage.Read(ctx, filter)
	if err != nil {
		return nil, err
	}

	return flights, nil
}

func (s *FlightService) CreateFlight(ctx context.Context, flightNumber string,
	origin string,
	destination string,
	date time.Time,
	status models.FlightStatus,
	aircraft string) (*models.Flight, error) {

	flight := models.NewFlight(flightNumber, origin, destination, date, status, aircraft)
	created, err := s.storage.Create(ctx, flight)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *FlightService) UpdateFlight(ctx context.Context, flight *models.Flight) (*models.Flight, error) {
	updated, err := s.storage.Update(ctx, flight)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *FlightService) DeleteFlight(ctx context.Context, id uuid.UUID) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *FlightService) CancelFlight(ctx context.Context, id uuid.UUID) error {
	flight, err := s.storage.ReadById(ctx, id)
	if err != nil {
		return err
	}
	if flight.Cancel() {
		_, err = s.storage.Update(ctx, flight)
		if err != nil {
			return err
		}
	}

	return nil
}
