package organizations

import (
	"errors"
)

var (
	errOrganizationIDRequired = errors.New("organization id required")
	errCommunityIDRequired    = errors.New("community id required")

	errAccountCreation        = errors.New("account creation")
	errAccountActivationEmail = errors.New("account activation email")

	errUploadLogo         = errors.New("upload logo")
	errUpdateOrganization = errors.New("update organization")
	errGeocode            = errors.New("geocode")
)
