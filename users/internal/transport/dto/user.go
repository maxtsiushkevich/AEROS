package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	Name     string    `json:"name" binding:"required,min=2,max=200,alphaunicode,excludesall=0123456789"`
	Email    string    `json:"email" binding:"required,email,lowercase,max=255"`
	Birthday time.Time `json:"birthday" binging:"required,pastdate"`
	Password string    `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Birthday    time.Time  `json:"birthday"`
	Email       string     `json:"email"`
	IsActivated bool       `json:"is_activated"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TokensResponse struct {
	AccessToken string `json:"access_token"`
}
