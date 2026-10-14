package dto

import (
	"github.com/google/uuid"
)

type CargoManifestResponse struct {
	ID              uuid.UUID                       `json:"id"`
	FlightID        uuid.UUID                       `json:"flight_id"`
	MaxWeightKg     int                             `json:"max_weight_kg"`
	CurrentWeightKg int                             `json:"current_weight_kg"`
	Items           map[uuid.UUID]CargoItemResponse `json:"items"`
}

type CargoItemResponse struct {
	ID          uuid.UUID  `json:"id"`
	CargoType   string     `json:"cargo_type"`
	WeightKg    int        `json:"weight_kg"`
	PassengerID *uuid.UUID `json:"passenger_id,omitempty"`
	Description string     `json:"description"`
}

type AddCargoItemRequest struct {
	CargoType   string     `json:"cargo_type" validate:"required,oneof=Luggage Cargo Mail Equipment Dangerous"`
	WeightKg    int        `json:"weight_kg" validate:"required,gt=0"`
	PassengerID *uuid.UUID `json:"passenger_id,omitempty" validate:"omitempty"`
	Description string     `json:"description" validate:"required,min=5"`
}
