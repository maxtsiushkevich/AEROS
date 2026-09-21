package models

import (
	"time"
	"users/internal/domain/errors"
	"users/internal/domain/rules"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID
	Name        string
	Birthday    time.Time
	Email       string
	IsActivated bool
	ActivatedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewUser(name string, email string, birthday time.Time) (*User, error) {
	if name == "" {
		return nil, errors.ErrInvalidName
	}

	if !rules.IsValidEmail(email) {
		return nil, errors.ErrInvalidEmail
	}

	err := rules.IsValidBirthday(birthday)
	if err != nil {
		return nil, err
	}

	if !rules.IsAdult(birthday) {
		return nil, errors.ErrNotAdult
	}

	return &User{
		Id:          uuid.New(),
		Name:        name,
		Email:       email,
		Birthday:    birthday,
		IsActivated: false,
		ActivatedAt: nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (u *User) Activate() error {
	if u.IsActivated {
		return errors.ErrUserAlreadyActivated
	}

	u.IsActivated = true

	now := time.Now()
	u.ActivatedAt = &now

	return nil
}
