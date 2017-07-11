package memberships

import (
	"errors"
)

var (
	errMembershipIDRequired = errors.New("membership id required")

	errAccountCreationFailed = errors.New("account creation failed")
	errAccountActivationEmail = errors.New("account activation email")
)
