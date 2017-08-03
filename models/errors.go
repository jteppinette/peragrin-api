package models

import (
	"errors"
)

var (
	errGeocodeNotFound = errors.New("geocode not found")
	errAccountNotFound    = errors.New("account not found")
	errInvalidCredentials = errors.New("invalid credentials")
)
