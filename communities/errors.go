package communities

import (
	"errors"
)

var (
	errCommunityIDRequired = errors.New("community id required")
	errCreateOrganization  = errors.New("create organization")

	errAuthenticationRequired = errors.New("authentication required")
	errSuperUserRequired = errors.New("super user required")

	errAccountActivationEmail = errors.New("account activation email")
)
