package errors

import (
	"errors"
)

var (
	ErrFlightNotFound      = errors.New("flight not found")
	ErrIncorrectFlightTime = errors.New("flight time in the past")
	ErrInvalidFlightStatus = errors.New("invalid flight status")
)

var (
	ErrIncorrectCargoWeigh = errors.New("incorrect cargo weight")
	ErrLuggageWithoutOwner = errors.New("luggage without owner")
	ErrMaxWeightExceeded   = errors.New("max weight exceeded")
)
