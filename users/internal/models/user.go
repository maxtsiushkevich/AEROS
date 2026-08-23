package models

import (
	"github.com/google/uuid"
)

type User struct {
	id    uuid.UUID
	name  string
	email string
}
