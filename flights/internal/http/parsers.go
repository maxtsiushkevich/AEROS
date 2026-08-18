package http

import (
	"net/url"
	"time"
)

// ParseGetFlightsQuery converts URL query parameters to GetFlightsRequestQuery
func ParseGetFlightsQuery(queryParams url.Values) (*GetFlightsRequestQuery, error) {
	dateFrom, err := parseDate(queryParams.Get("date_from"))
	if err != nil {
		return nil, err
	}

	dateTo, err := parseDate(queryParams.Get("date_to"))
	if err != nil {
		return nil, err
	}

	return &GetFlightsRequestQuery{
		FlightNumber: queryParams.Get("flight_number"),
		Origin:       queryParams.Get("origin"),
		Destination:  queryParams.Get("destination"),
		Status:       queryParams.Get("status"),
		DateFrom:     dateFrom,
		DateTo:       dateTo,
	}, nil
}

// parseDate parses RFC3339 date string, returns nil if empty
func parseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
