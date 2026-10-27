package events

type Event string

const (
	UserCreated  Event = "UserCreated"
	EmailChanged Event = "EmailChanged"
)
