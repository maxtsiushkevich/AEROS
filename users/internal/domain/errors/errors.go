package errors

import "errors"

var (
	ErrUserAlreadyActivated = errors.New("user is already activated")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidName          = errors.New("invalid name")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrBirthdayInvalid      = errors.New("invalid birthday date")
	ErrNotAdult             = errors.New("not adult")
)
