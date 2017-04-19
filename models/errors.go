package models

import (
	"errors"
)

var (
	errGeo             = errors.New("geo")
	errAddressRequired = errors.New("address requird")
)
