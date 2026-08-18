package storage

import (
	"time"
)

type FlightFilter struct {
	FlightNumber string
	Origin       string
	Destination  string
	Status       string
	DateFrom     *time.Time
	DateTo       *time.Time
}
