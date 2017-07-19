package accounts

import "errors"

var (
	errAccountIDRequired = errors.New("account id required")

	errAccountNotFound = errors.New("account not found")

	errCreateOrganization = errors.New("create organization")
	errGeocode            = errors.New("geocode")
)
