package communities

import (
	"errors"
)

var (
	errCommunityIDRequired = errors.New("community id required")
	errCreateOrganization  = errors.New("create organization")

	errAccountActivationEmail = errors.New("account activation email")
)
