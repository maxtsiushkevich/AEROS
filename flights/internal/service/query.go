package service

import (
	"flights/internal/domain"
	"time"

	_ "github.com/go-playground/validator/v10"
)

type FlightQuery struct {
	FlightNumber string
	Origin       string
	Destination  string
	Status       domain.Status
	DateFrom     *time.Time
	DateTo       *time.Time
}
