package service

import (
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

func (s *FlightService) GetFlights(query *models.FlightQuery) ([]models.Flight, error) {
	// Could add filtering rules, business rules checking, etc.
	flights, err := s.storage.Read(query)
	if err != nil {
		return nil, err
	}

	return flights, nil
}

func (s *FlightService) CreateFlight(flight *models.Flight) error {
	// check if aircraft exists, validate route, etc.
	err := s.storage.Create(flight)
	if err != nil {
		return err
	}

	return nil
}
