package rules

import (
	"net/mail"
	"time"
	"users/internal/domain/errors"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidBirthday(date time.Time) error {
	now := time.Now()

	if date.After(now) {
		return errors.ErrBirthdayInvalid
	}

	if now.Year()-date.Year() > 120 {
		return errors.ErrBirthdayInvalid
	}

	return nil
}

func IsAdult(birthday time.Time) bool {
	adultDate := birthday.AddDate(18, 0, 0)
	return !time.Now().Before(adultDate)
}
