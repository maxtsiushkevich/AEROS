package postgres

import "flights/internal/domain/models"

func FlightToDomain(f *Flight) *models.Flight {
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

func FlightsToDomain(dbFlights []Flight) []models.Flight {
	flights := make([]models.Flight, 0, len(dbFlights))
	for _, dbFlight := range dbFlights {
		if flight := FlightToDomain(&dbFlight); flight != nil {
			flights = append(flights, *flight)
		}
	}
	return flights
}

func FlightFromDomain(f *models.Flight) *Flight {
	if f == nil {
		return nil
	}

	dbFlight := &Flight{}
	FlightToDb(f, dbFlight)
	return dbFlight
}

func FlightToDb(f *models.Flight, dbFlight *Flight) {
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

func (i CargoItem) ToDomain() models.CargoItem {
	return models.CargoItem{
		ID:          i.ID,
		ManifestID:  i.ManifestID,
		CargoType:   models.CargoType(i.CargoType),
		WeightKg:    i.WeightKg,
		PassengerID: i.PassengerID,
		Description: i.Description,
	}
}

func CargoItemFromDomain(d models.CargoItem) CargoItem {
	return CargoItem{
		ID:          d.ID,
		ManifestID:  d.ManifestID,
		CargoType:   models.CargoType(d.CargoType),
		WeightKg:    d.WeightKg,
		PassengerID: d.PassengerID,
		Description: d.Description,
	}
}

func (m CargoManifest) ToDomain() *models.CargoManifest {
	items := make([]models.CargoItem, 0, len(m.CargoItems))
	for _, it := range m.CargoItems {
		items = append(items, it.ToDomain())
	}

	return models.RestoreCargoManifest(
		m.ID,
		m.FlightID,
		m.MaxWeightKg,
		m.CurrentWeightKg,
		items,
	)
}

func CargoManifestFromDomain(d *models.CargoManifest) CargoManifest {
	domainItems := d.Items()
	items := make([]CargoItem, 0, len(domainItems))
	for _, it := range domainItems {
		items = append(items, CargoItemFromDomain(it))
	}

	return CargoManifest{
		Header:          Header{ID: d.ID},
		FlightID:        d.FlightID,
		MaxWeightKg:     d.MaxWeightKg,
		CurrentWeightKg: d.CurrentWeightKg,
		CargoItems:      items,
	}
}
