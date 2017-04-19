package organizations

import (
	"errors"
)

var (
	errListOrganizations = errors.New("list organizations")
	errGetOrganization   = errors.New("get organization")

	errOrganizationIDRequired = errors.New("organization id required")
)
