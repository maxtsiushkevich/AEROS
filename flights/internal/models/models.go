package models

import (
	"time"

	"gorm.io/gorm"
)

type Status int

const (
	Scheduled Status = iota
	CheckIn
	Boarding
	Delayed
	Departed
	Arrived
	Cancelled
	Redirected
)

type Flight struct {
	gorm.Model
	FlightNumber string `gorm:"size:8"`
	Origin       string `gorm:"size:3"`
	Destination  string `gorm:"size:3"`
	Date         time.Time
	Status       Status
	Aircraft     string
}
