package service

import (
	"context"
	"flights/internal/domain/models"
	"time"
)

type FlightService interface {
	CanReschedule(ctx context.Context, flight *models.Flight, newDate time.Time) (bool, error)
	Reschedule(ctx context.Context, flight *models.Flight, newDate time.Time) error
}
