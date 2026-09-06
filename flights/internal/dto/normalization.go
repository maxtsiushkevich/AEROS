package dto

import "strings"

func (r *GetFlightsRequestQuery) Normalize() {
	r.FlightNumber = trimUpper(r.FlightNumber)
	r.Origin = trimUpper(r.Origin)
	r.Destination = trimUpper(r.Destination)

	r.Status = trim(r.Status)
}

func (r *CreateFlightRequest) Normalize() {
	r.FlightNumber = strings.ToUpper(strings.TrimSpace(r.FlightNumber))
	r.Origin = strings.ToUpper(strings.TrimSpace(r.Origin))
	r.Destination = strings.ToUpper(strings.TrimSpace(r.Destination))
	r.Status = strings.TrimSpace(r.Status)
	if r.Status == "" {
		r.Status = "Scheduled"
	}
	r.Aircraft = strings.TrimSpace(r.Aircraft)
}

func (r *PatchFlightRequest) Normalize() {
	r.FlightNumber = trimUpper(r.FlightNumber)
	r.Origin = trimUpper(r.Origin)
	r.Destination = trimUpper(r.Destination)
	r.Aircraft = trim(r.Aircraft)
	r.Status = trim(r.Status)
}

func trim(s *string) *string {
	if s == nil {
		return nil
	}
	res := strings.TrimSpace(*s)
	return &res
}

func trimUpper(s *string) *string {
	if s == nil {
		return nil
	}
	res := strings.ToUpper(strings.TrimSpace(*s))
	return &res
}
