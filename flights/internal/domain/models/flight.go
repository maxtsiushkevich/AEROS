package models

import (
	"flights/internal/domain/errors"
	"time"

	"github.com/google/uuid"
)

type FlightStatus string

const (
	Scheduled  FlightStatus = "Scheduled"
	CheckIn    FlightStatus = "CheckIn"
	Boarding   FlightStatus = "Boarding"
	Delayed    FlightStatus = "Delayed"
	Departed   FlightStatus = "Departed"
	Arrived    FlightStatus = "Arrived"
	Cancelled  FlightStatus = "Cancelled"
	Redirected FlightStatus = "Redirected"
)

type Flight struct {
	Id           uuid.UUID
	FlightNumber string
	Origin       string
	Destination  string
	Date         time.Time
	Status       FlightStatus
	Aircraft     string
}

func NewFlight(flightNumber string,
	origin string,
	destination string,
	date time.Time,
	status FlightStatus,
	aircraft string) *Flight {
	return &Flight{
		Id:           uuid.New(),
		FlightNumber: flightNumber,
		Origin:       origin,
		Destination:  destination,
		Date:         date,
		Status:       status,
		Aircraft:     aircraft,
	}
}

func (f *Flight) Cancel() bool {
	if f.Status != Arrived && f.Status != Departed && f.Status != Cancelled {
		f.Status = Cancelled
		return true
	}
	return false
}

func (f *Flight) Reschedule(newDate time.Time) error {
	if !newDate.After(time.Now()) {
		return errors.ErrIncorrectFlightTime
	}

	if newDate.After(f.Date) {
		f.Status = Delayed
	}

	f.Date = newDate
	return nil
}

func (f *Flight) ChangeStatus(status string) error {
	if !IsValidFlightStatus(status) {
		return errors.ErrInvalidFlightStatus
	}

	f.Status = FlightStatus(status)
	return nil
}
