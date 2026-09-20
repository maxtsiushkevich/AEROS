package models

import (
	"time"
)

type User struct {
	Name         string
	Birthday     time.Time
	Email        string
	MemberNumber string
	IsActivated  bool
	ActivatedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateUser(name string, email string, birthday time.Time) (*User, error) {
	return nil, nil
}

func (u *User) Activate() error {
	return nil
}

func (u *User) Register() error {
	return nil
}
