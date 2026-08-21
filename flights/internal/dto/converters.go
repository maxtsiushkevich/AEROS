package dto

import (
	"flights/internal/models"
)

func (r *CreateFlightRequest) ToServiceFlight() *models.Flight {
	return &models.Flight{
		FlightNumber: *r.FlightNumber,
		Origin:       *(r.Origin),
		Destination:  *r.Destination,
		Date:         *r.Date,
		Status:       models.Status(*r.Status),
		Aircraft:     *r.Aircraft,
	}
}

func FlightToResponse(flight *models.Flight) FlightResponse {
	return FlightResponse{
		ID:           flight.ID,
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

func (r *GetFlightsRequestQuery) ToServiceQuery() *models.FlightQuery {
	return &models.FlightQuery{
		FlightNumber: r.FlightNumber,
		Origin:       r.Origin,
		Destination:  r.Destination,
		Status:       r.Status,
		DateFrom:     r.DateFrom,
		DateTo:       r.DateTo,
	}
}

func (r *PatchFlightRequest) ToFlightUpdate() *models.FlightUpdate {
	update := &models.FlightUpdate{
		ID:           r.ID,
		FlightNumber: r.FlightNumber,
		Origin:       r.Origin,
		Destination:  r.Destination,
		Date:         r.Date,
		Aircraft:     r.Aircraft,
	}

	if r.Status != nil {
		status := models.Status(*r.Status)
		update.Status = &status
	}

	return update
}
