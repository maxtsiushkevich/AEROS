package dto

import (
	"flights/internal/domain/models"

	"github.com/google/uuid"
)

func (r *GetFlightsRequestQuery) ToFlightFilter() *models.FlightFilter {
	return &models.FlightFilter{
		FlightNumber: r.FlightNumber,
		Origin:       r.Origin,
		Destination:  r.Destination,
		Status:       r.Status,
		DateFrom:     r.DateFrom,
		DateTo:       r.DateTo,
		Limit:        r.Limit,
		Page:         r.Page,
	}
}

func FlightToResponse(flight *models.Flight) FlightResponse {
	return FlightResponse{
		ID:           flight.Id,
		FlightNumber: flight.FlightNumber,
		Origin:       flight.Origin,
		Destination:  flight.Destination,
		Date:         flight.Date,
		Status:       string(flight.Status),
		Aircraft:     flight.Aircraft,
	}
}

func FlightsToResponses(flights []models.Flight) []FlightResponse {
	responses := make([]FlightResponse, len(flights))
	for i, flight := range flights {
		responses[i] = FlightToResponse(&flight)
	}
	return responses
}

func CargoManifestToResponse(cm *models.CargoManifest) CargoManifestResponse {
	domainItems := cm.Items()

	items := make(map[uuid.UUID]CargoItemResponse, len(domainItems))
	for _, item := range domainItems {
		items[item.ID] = CargoItemToResponse(&item)
	}

	return CargoManifestResponse{
		ID:              cm.ID,
		FlightID:        cm.FlightID,
		MaxWeightKg:     cm.MaxWeightKg,
		CurrentWeightKg: cm.CurrentWeightKg,
		Items:           items,
	}
}

func CargoItemToResponse(ci *models.CargoItem) CargoItemResponse {
	return CargoItemResponse{
		ID:          ci.ID,
		CargoType:   string(ci.CargoType),
		WeightKg:    ci.WeightKg,
		PassengerID: ci.PassengerID,
		Description: ci.Description,
	}
}
