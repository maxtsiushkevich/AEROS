package service

import (
	"flights/internal/domain"
	"flights/internal/storage"
)

type FlightService struct {
	storage storage.Storage
}

func CreateFlightService(storage storage.Storage) *FlightService {
	return &FlightService{
		storage: storage,
	}
}

func (s *FlightService) GetFlights(flightQuery FlightQuery) ([]storage.Flight, error) {
	filter := storage.FlightFilter{
		FlightNumber: flightQuery.FlightNumber,
		Origin:       flightQuery.Origin,
		Destination:  flightQuery.Destination,
		Status:       string(flightQuery.Status),
		DateFrom:     flightQuery.DateFrom,
		DateTo:       flightQuery.DateTo,
	}

	flights, err := s.storage.Read(filter)
	if err != nil {
		return nil, err
	}
	return flights, nil
}

func (s *FlightService) CreateFlight(flight *CreateFlightModel) error {
	err := s.storage.Create(&storage.Flight{
		FlightNumber: flight.FlightNumber,
		Origin:       flight.Origin,
		Destination:  flight.Destination,
		Date:         flight.Date,
		Status:       domain.Status(flight.Status),
		Aircraft:     flight.Aircraft,
	})

	if err != nil {
		return err
	}
	return nil
}
