package organizations

import (
	"errors"
)

var (
	errCreateOrganization = errors.New("create organization")
	errUpdateOrganization = errors.New("update organization")
	errListOrganizations  = errors.New("list organizations")
	errGetOrganization    = errors.New("get organization")

	errOrganizationIDRequired = errors.New("organization id required")
	errGeocode                = errors.New("geocode")
)
