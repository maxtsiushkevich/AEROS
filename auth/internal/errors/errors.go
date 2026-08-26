package errors

import "fmt"

type AuthError struct {
	Code    string
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewAuthError(code, message string) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
	}
}

func UserNotFoundError(email string) *AuthError {
	return &AuthError{
		Code:    "USER_NOT_FOUND",
		Message: fmt.Sprintf("user with email %s not found", email),
	}
}

func UserAlreadyExistsError(email string) *AuthError {
	return &AuthError{
		Code:    "USER_ALREADY_EXISTS",
		Message: fmt.Sprintf("user with email %s already exists", email),
	}
}

func InvalidPasswordError(reason string) *AuthError {
	return &AuthError{
		Code:    "INVALID_PASSWORD",
		Message: reason,
	}
}

func CreateTokenError(reason string) *AuthError {
	return &AuthError{
		Code:    "TOKEN_ERROR",
		Message: reason,
	}
}
