package models

import (
	"flights/internal/domain/errors"

	"github.com/google/uuid"
)

type CargoManifest struct {
	ID              uuid.UUID
	FlightID        uuid.UUID
	MaxWeightKg     int
	CurrentWeightKg int
	items           map[uuid.UUID]CargoItem
}

func NewCargoManifest(flightID uuid.UUID, maxWeightKg int) (*CargoManifest, error) {
	if maxWeightKg <= 0 {
		return nil, errors.ErrIncorrectCargoWeigh
	}

	return &CargoManifest{
		ID:              uuid.New(),
		FlightID:        flightID,
		MaxWeightKg:     maxWeightKg,
		CurrentWeightKg: 0,
		items:           map[uuid.UUID]CargoItem{},
	}, nil
}

func (m *CargoManifest) Items() []CargoItem {
	res := make([]CargoItem, 0, len(m.items))
	for _, it := range m.items {
		res = append(res, it)
	}
	return res
}

func RestoreCargoManifest(
	id, flightID uuid.UUID,
	maxWeightKg, currentWeightKg int,
	items []CargoItem,
) *CargoManifest {
	m := &CargoManifest{
		ID:              id,
		FlightID:        flightID,
		MaxWeightKg:     maxWeightKg,
		CurrentWeightKg: currentWeightKg,
		items:           make(map[uuid.UUID]CargoItem, len(items)),
	}
	for _, it := range items {
		m.items[it.ID] = it
	}
	return m
}

func (m *CargoManifest) AddCargoItem(item CargoItem) error {
	if m.CurrentWeightKg+item.WeightKg > m.MaxWeightKg {
		return errors.ErrMaxWeightExceeded
	}

	m.items[item.ID] = item
	m.CurrentWeightKg += item.WeightKg
	return nil
}

func (m *CargoManifest) DeleteCargoItem(itemID uuid.UUID) {
	if item, exists := m.items[itemID]; exists {
		m.CurrentWeightKg -= item.WeightKg
		if m.CurrentWeightKg < 0 {
			m.CurrentWeightKg = 0
		}
	}
	delete(m.items, itemID)
}

func (m *CargoManifest) CargoItems() []CargoItem {
	items := make([]CargoItem, 0, len(m.items))
	for _, item := range m.items {
		items = append(items, item)
	}
	return items
}
