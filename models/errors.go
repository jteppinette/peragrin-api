package models

import (
	"errors"
)

var (
	errAddressRequired = errors.New("address requird")
	errGeocodeNotFound = errors.New("geocode not found")

	errAccountNotFound    = errors.New("account not found")
	errInvalidCredentials = errors.New("invalid credentials")
)
