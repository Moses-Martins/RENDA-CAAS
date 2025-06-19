package utils

import (
	"errors"
	"unicode"
)

// ValidatePassword checks if the password is at least 8 characters long,
// contains both letters and numbers, and does not contain symbols.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	hasLetter := false
	hasNumber := false
	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsNumber(c):
			hasNumber = true
		case !unicode.IsLetter(c) && !unicode.IsNumber(c):
			return errors.New("password must not contain symbols")
		}
	}
	if !hasLetter || !hasNumber {
		return errors.New("password must contain both letters and numbers")
	}
	return nil
}
