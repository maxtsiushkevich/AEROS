package dto

import (
	"time"
)

type User struct {
	Name     string    `json:"name" validate:"required,min=2,max=200,alphaunicode,excludesall=0123456789"`
	Email    string    `json:"email" validate:"required,email,lowercase,max=255"`
	Birthday time.Time `json:"birthday" validate:"required,pastdate"`
}
