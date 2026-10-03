package models

import (
	"flights/internal/domain/errors"
	"time"

	"github.com/google/uuid"
)

type CargoItem struct {
	Id          uuid.UUID
	CargoType   CargoType
	WeightKg    int
	PassengerId *uuid.UUID
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
	cargo_type CargoType,
	weight_kg int,
	passengerId *uuid.UUID,
	description string,
	packed_at time.Time) (*CargoItem, error) {

	if weight_kg <= 0 {
		return nil, errors.ErrIncorrectCargoWeigh
	}

	if passengerId == nil && cargo_type == Luggage {
		return nil, errors.ErrLuggageWithoutOwner
	}

	return &CargoItem{
		Id:          uuid.New(),
		CargoType:   cargo_type,
		WeightKg:    weight_kg,
		PassengerId: passengerId,
		Description: description,
		PackedAt:    packed_at,
	}, nil
}
