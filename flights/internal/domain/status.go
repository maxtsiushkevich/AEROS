package domain

type Status string

const (
	Scheduled  Status = "Scheduled"
	CheckIn    Status = "CheckIn"
	Boarding   Status = "Boarding"
	Delayed    Status = "Delayed"
	Departed   Status = "Departed"
	Arrived    Status = "Arrived"
	Cancelled  Status = "Cancelled"
	Redirected Status = "Redirected"
)
