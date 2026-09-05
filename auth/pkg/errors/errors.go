package errors

import (
	"errors"
	"fmt"
)

type AuthError struct {
	Code    string
	Message string
}

var SamePasswordError = errors.New("New password cannot be same as old password")
var UserNotFoundError = errors.New("user with this email not found")
var UserAlreadyExistsError = errors.New("user with this email already exists")
var InvalidPasswordError = errors.New("incorrect password")

type CreateTokenError error

func (e *AuthError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewCreateTokenError(reason string) CreateTokenError {
	return errors.New(reason)
}

func NewAuthError(code, message string) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
	}
}
