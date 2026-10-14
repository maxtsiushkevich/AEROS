package models

import (
	"flights/internal/domain/errors"

	"github.com/google/uuid"
)

type CargoItem struct {
	ID          uuid.UUID
	ManifestID  uuid.UUID
	CargoType   CargoType
	WeightKg    int
	PassengerID *uuid.UUID
	Description string
}

type CargoType string

const (
	Luggage      CargoType = "Luggage"   // Багаж пассажиров
	Cargo        CargoType = "Cargo"     // Коммерческий груз
	Mail         CargoType = "Mail"      // Почта
	Equipment    CargoType = "Equipment" // Оборудование
	DangerousGds CargoType = "Dangerous" // Опасные грузы
)

func NewCargoItem(
	manifestID uuid.UUID,
	cargoType CargoType,
	weightKg int,
	passengerID *uuid.UUID,
	description string,
) (*CargoItem, error) {
	if weightKg <= 0 {
		return nil, errors.ErrIncorrectCargoWeigh
	}

	if passengerID == nil && cargoType == Luggage {
		return nil, errors.ErrLuggageWithoutOwner
	}

	return &CargoItem{
		ID:          uuid.New(),
		ManifestID:  manifestID,
		CargoType:   cargoType,
		WeightKg:    weightKg,
		PassengerID: passengerID,
		Description: description,
	}, nil
}
