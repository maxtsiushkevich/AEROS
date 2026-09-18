package utils

import (
	"flights/internal/dto"
	"net/url"
	"strconv"
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

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	return &dto.GetFlightsRequestQuery{
		FlightNumber: getStringOrNil(queryParams, "flight_number"),
		Origin:       getStringOrNil(queryParams, "origin"),
		Destination:  getStringOrNil(queryParams, "destination"),
		Status:       getStringOrNil(queryParams, "status"),
		DateFrom:     dateFrom,
		DateTo:       dateTo,
		Limit:        limit,
		Page:         page,
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
