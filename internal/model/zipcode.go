package model

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidZipcode  = errors.New("invalid zipcode")
	ErrZipcodeNotFound = errors.New("zipcode not found")
	ErrWeatherFailure  = errors.New("weather api failure")
)

var zipcodeRegex = regexp.MustCompile(`^\d{8}$`)

func ValidateZipcode(zipcode string) error {
	if !zipcodeRegex.MatchString(zipcode) {
		return ErrInvalidZipcode
	}

	return nil
}
