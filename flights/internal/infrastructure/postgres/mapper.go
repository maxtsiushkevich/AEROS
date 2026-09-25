package postgres

import "flights/internal/domain/models"

func ToDomain(f *Flight) *models.Flight {
	if f == nil {
		return nil
	}

	return &models.Flight{
		Id:           f.ID,
		FlightNumber: f.FlightNumber,
		Origin:       f.Origin,
		Destination:  f.Destination,
		Date:         f.Date,
		Status:       f.Status,
		Aircraft:     f.Aircraft,
	}
}

func DomainToDB(f *models.Flight, dbFlight *Flight) {
	if f == nil || dbFlight == nil {
		return
	}

	dbFlight.ID = f.Id
	if f.FlightNumber != "" {
		dbFlight.FlightNumber = f.FlightNumber
	}
	if f.Origin != "" {
		dbFlight.Origin = f.Origin
	}
	if f.Destination != "" {
		dbFlight.Destination = f.Destination
	}
	if !f.Date.IsZero() {
		dbFlight.Date = f.Date
	}
	if f.Status != "" {
		dbFlight.Status = f.Status
	}
	if f.Aircraft != "" {
		dbFlight.Aircraft = f.Aircraft
	}
}
