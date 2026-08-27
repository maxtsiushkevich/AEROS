package rbac

type Role string

const (
	Admin            Role = "admin"
	User             Role = "user"
	SalesAgent       Role = "salesAgent"
	CustromerService Role = "customerService"
	Pilot            Role = "pilot"
	Dispatcher       Role = " dispatcher"
)
