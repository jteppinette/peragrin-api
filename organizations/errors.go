package organizations

import (
	"errors"
)

var (
	errOrganizationIDRequired = errors.New("organization id required")
	errAccountIDRequired      = errors.New("account id required")
	errCommunityIDRequired    = errors.New("community id required")

	errJoinCommunity = errors.New("join community")

	errAccountCreation        = errors.New("account creation")
	errAccountActivationEmail = errors.New("account activation email")

	errAuthenticationRequired = errors.New("authentication required")

	errUploadLogo         = errors.New("upload logo")
	errUpdateOrganization = errors.New("update organization")
)
