package service

import "time"

type CreateFlightModel struct {
	FlightNumber string
	Origin       string
	Destination  string
	Date         time.Time
	Status       string
	Aircraft     string
}
