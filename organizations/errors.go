package organizations

import (
	"errors"
)

var (
	errCreateOrganization = errors.New("create organization")
	errUpdateOrganization = errors.New("update organization")
	errListOrganizations  = errors.New("list organizations")
	errGetOrganization    = errors.New("get organization")

	errCreatePost = errors.New("create post")

	errOrganizationIDRequired = errors.New("organization id required")
	errGeocode                = errors.New("geocode")

	errAddOperator = errors.New("add operator")

	errAuthenticationRequired = errors.New("authentication required")
)
