package models

import (
	"flights/internal/domain/errors"
	"time"

	"github.com/google/uuid"
)

type CargoItem struct {
	ID          uuid.UUID
	ManifestID  uuid.UUID
	CargoType   CargoType
	WeightKg    int
	PassengerID *uuid.UUID
	Description string
	PackedAt    time.Time
}

type CargoType string

const (
	Luggage      CargoType = "LUGGAGE"   // Багаж пассажиров
	Cargo        CargoType = "CARGO"     // Коммерческий груз
	Mail         CargoType = "MAIL"      // Почта
	Equipment    CargoType = "EQUIPMENT" // Оборудование
	DangerousGds CargoType = "DANGEROUS" // Опасные грузы
)

func NewCargoItem(
	cargoType CargoType,
	weightKg int,
	passengerID *uuid.UUID,
	description string,
	packedAt time.Time,
) (*CargoItem, error) {
	if weightKg <= 0 {
		return nil, errors.ErrIncorrectCargoWeigh
	}

	if passengerID == nil && cargoType == Luggage {
		return nil, errors.ErrLuggageWithoutOwner
	}

	return &CargoItem{
		ID:          uuid.New(),
		CargoType:   cargoType,
		WeightKg:    weightKg,
		PassengerID: passengerID,
		Description: description,
		PackedAt:    packedAt,
	}, nil
}
