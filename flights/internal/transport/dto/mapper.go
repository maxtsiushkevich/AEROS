package dto

import "flights/internal/domain/models"

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
