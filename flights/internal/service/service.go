package service

import (
	"context"
	"flights/internal/models"
	"flights/internal/storage"
)

type FlightService struct {
	storage storage.FlightsStorage
}

func CreateFlightService(storage storage.FlightsStorage) *FlightService {
	return &FlightService{
		storage: storage,
	}
}

func (s *FlightService) GetFlights(ctx context.Context, query *models.FlightQuery) ([]models.Flight, error) {
	// Could add filtering rules, business rules checking, etc.
	flights, err := s.storage.Read(ctx, query)
	if err != nil {
		return nil, err
	}

	return flights, nil
}

func (s *FlightService) CreateFlight(ctx context.Context, flight *models.Flight) error {
	// check if aircraft exists, validate route, etc.
	err := s.storage.Create(ctx, flight)
	if err != nil {
		return err
	}

	return nil
}

func (s *FlightService) UpdateFlight(ctx context.Context, flight *models.FlightUpdate) (*models.Flight, error) {
	return &models.Flight{}, nil
}
