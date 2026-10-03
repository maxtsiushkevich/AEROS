package usecase

import (
	"context"
	"flights/internal/domain/models"
	"flights/internal/domain/repository"
	"flights/internal/domain/service"
	"time"

	"github.com/google/uuid"
)

type FlightUseCase struct {
	flightsStorage repository.FlightRepository
	cargoStorage   repository.CargoRepository
	service        service.FlightDomainService
}

func NewFlightUseCase(storage repository.FlightRepository, cargoStorage repository.CargoRepository) *FlightUseCase {
	return &FlightUseCase{
		flightsStorage: storage,
		cargoStorage:  cargoStorage,
	}
}

func (s *FlightUseCase) GetFlights(ctx context.Context, filter *models.FlightFilter) ([]models.Flight, error) {
	flights, err := s.flightsStorage.Read(ctx, filter)
	if err != nil {
		return nil, err
	}

	return flights, nil
}

func (s *FlightUseCase) CreateFlight(ctx context.Context, flightNumber string,
	origin string,
	destination string,
	date time.Time,
	status models.FlightStatus,
	aircraft string) (*models.Flight, error) {

	flight := models.NewFlight(flightNumber, origin, destination, date, status, aircraft)
	created, err := s.flightsStorage.Create(ctx, flight)
	if err != nil {
		return nil, err
	}

	cargo, err := models.NewCargoManifest(flight.Id, 20000)
	if err != nil {
		return nil, err
	}

	err = s.cargoStorage.Create(ctx, cargo)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *FlightUseCase) UpdateFlight(ctx context.Context, flight *models.Flight) (*models.Flight, error) {
	updated, err := s.flightsStorage.Update(ctx, flight)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *FlightUseCase) DeleteFlight(ctx context.Context, id uuid.UUID) error {
	err := s.flightsStorage.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = s.cargoStorage.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *FlightUseCase) CancelFlight(ctx context.Context, id uuid.UUID) error {
	flight, err := s.flightsStorage.ReadById(ctx, id)
	if err != nil {
		return err
	}
	if flight.Cancel() {
		_, err = s.flightsStorage.Update(ctx, flight)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FlightUseCase) RescheduleFlight(ctx context.Context, id uuid.UUID, newDate time.Time) error {
	f, err := s.flightsStorage.ReadById(ctx, id)
	if err != nil {
		return err
	}

	err = s.service.Reschedule(ctx, f, newDate)
	if err != nil {
		return err
	}

	_, err = s.flightsStorage.Update(ctx, f)
	if err != nil {
		return err
	}

	return nil
}

func (s *FlightUseCase) ChangeStatus(ctx context.Context, id uuid.UUID, status string) error {
	f, err := s.flightsStorage.ReadById(ctx, id)
	if err != nil {
		return err
	}

	err = f.ChangeStatus(status)
	if err != nil {
		return err
	}

	_, err = s.flightsStorage.Update(ctx, f)
	if err != nil {
		return err
	}

	return nil
}
