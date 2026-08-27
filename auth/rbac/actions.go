package rbac

type Action string

const (
	Read      Action = "read"
	Write     Action = "write"
	ReadWrite Action = "readwrite"
	Delete    Action = "delete"
)
