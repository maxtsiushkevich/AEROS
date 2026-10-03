package service

import (
	"context"
	"flights/internal/domain/errors"
	"flights/internal/domain/models"
	"flights/internal/domain/repository"
	"time"
)

type FlightDomainService struct {
	repo repository.FlightRepository
}

func NewFlightDomainService(repo repository.FlightRepository) *FlightDomainService {
	return &FlightDomainService{repo: repo}
}

func (s *FlightDomainService) CanReschedule(ctx context.Context, flight *models.Flight, newDate time.Time) (bool, error) {
	if !newDate.After(time.Now()) {
		return false, errors.ErrIncorrectFlightTime
	}

	dateFrom := flight.Date.Add(-5 * time.Hour)
	dateTo := flight.Date.Add(5 * time.Hour)

	allFlights, err := s.repo.Read(ctx, &models.FlightFilter{
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
	})

	if err != nil {
		return false, err
	}

	for _, candidate := range allFlights {
		if candidate.Id == flight.Id {
			continue
		}
		if candidate.Aircraft != flight.Aircraft {
			continue
		}
		if candidate.Date.Equal(newDate) {
			return false, errors.ErrFlightScheduleConflict
		}
	}

	return true, nil
}

func (s *FlightDomainService) Reschedule(ctx context.Context, flight *models.Flight, newDate time.Time) error {

	canReschedule, err := s.CanReschedule(ctx, flight, newDate)
	if err != nil {
		return err
	}
	if !canReschedule {
		return errors.ErrFlightScheduleConflict
	}

	if err = flight.Reschedule(newDate, time.Now()); err != nil {
		return err
	}

	return err
}
