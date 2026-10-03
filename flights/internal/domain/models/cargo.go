package models

import (
	"flights/internal/domain/errors"

	"github.com/google/uuid"
)

type CargoManifest struct {
	flightId        uuid.UUID
	maxWeightKg     int
	currentWeightKg int
	cargoItems      map[uuid.UUID]CargoItem
}

func NewCargoManifest(fligh_id uuid.UUID, max_weight_kg int) (*CargoManifest, error) {
	if max_weight_kg <= 0 {
		return nil, errors.ErrIncorrectCargoWeigh
	}
	return &CargoManifest{
		flightId:        fligh_id,
		maxWeightKg:     max_weight_kg,
		currentWeightKg: 0,
		cargoItems:      map[uuid.UUID]CargoItem{},
	}, nil
}

func (m *CargoManifest) AddCargoItem(item CargoItem) error {
	if m.currentWeightKg+item.WeightKg > m.maxWeightKg {
		return errors.ErrMaxWeightExceeded
	}

	m.cargoItems[item.Id] = item
	m.currentWeightKg += item.WeightKg
	return nil
}

func (m *CargoManifest) DeleteCargoItem(item_id uuid.UUID) {
	if item, exists := m.cargoItems[item_id]; exists {
		m.currentWeightKg -= item.WeightKg
	}
	delete(m.cargoItems, item_id)
}

func (m *CargoManifest) FlightId() uuid.UUID { return m.flightId }

func (m *CargoManifest) MaxWeightKg() int { return m.maxWeightKg }

func (m *CargoManifest) CurrentWeightKg() int { return m.currentWeightKg }

func (m *CargoManifest) CargoItems() []CargoItem {
	items := make([]CargoItem, 0, len(m.cargoItems))
	for _, item := range m.cargoItems {
		items = append(items, item)
	}
	return items
}
