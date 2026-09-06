package utils

import (
	"flights/internal/dto"
	"net/url"
	"time"
)

func getStringOrNil(queryParams url.Values, key string) *string {
	value := queryParams.Get(key)
	if value == "" {
		return nil
	}
	return &value
}

// ParseGetFlightsQuery converts URL query parameters to GetFlightsRequestQuery
func ParseGetFlightsQuery(queryParams url.Values) (*dto.GetFlightsRequestQuery, error) {
	dateFrom, err := parseDate(queryParams.Get("date_from"))
	if err != nil {
		return nil, err
	}

	dateTo, err := parseDate(queryParams.Get("date_to"))
	if err != nil {
		return nil, err
	}

	return &dto.GetFlightsRequestQuery{
		FlightNumber: getStringOrNil(queryParams, "flight_number"),
		Origin:       getStringOrNil(queryParams, "origin"),
		Destination:  getStringOrNil(queryParams, "destination"),
		Status:       getStringOrNil(queryParams, "status"),
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
