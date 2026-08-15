package service

import (
	"flights/internal/models"
	"flights/internal/storage"
	"fmt"
)

type FlightService struct {
	storage *storage.Storage
}

func CreateFlightService(storage *storage.Storage) *FlightService {
	return &FlightService{
		storage: storage,
	}
}

func (s *FlightService) GetFlightsByNumber(flightNumber string) (*models.Flight, error) {
	res, err := s.storage.ReadByFlightNumber(flightNumber)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	return res, nil
}
