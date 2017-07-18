package memberships

import (
	"errors"
)

var (
	errMembershipIDRequired = errors.New("membership id required")

	errAccountCreation        = errors.New("account creation")
	errAccountActivationEmail = errors.New("account activation email")
)
