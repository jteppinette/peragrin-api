package organizations

import (
	"errors"
)

var (
	errListOrganizations  = errors.New("list organizations")
	errGetOrganization    = errors.New("get organization")
	errEnableOrganization = errors.New("enable organization")

	errOrganizationIDRequired = errors.New("organization id required")
)
