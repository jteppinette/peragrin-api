package models

import (
	"errors"
)

var (
	errAddressRequired = errors.New("address requird")
	errGeocodeNotFound = errors.New("geocode not found")
)
