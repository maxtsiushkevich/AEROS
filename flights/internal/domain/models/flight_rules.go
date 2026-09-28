package models

var validStatuses = map[FlightStatus]bool{
	Scheduled:  true,
	CheckIn:    true,
	Boarding:   true,
	Delayed:    true,
	Departed:   true,
	Arrived:    true,
	Cancelled:  true,
	Redirected: true,
}

func IsValidFlightStatus(statusStr string) bool {
	_, exists := validStatuses[FlightStatus(statusStr)]
	return exists
}
