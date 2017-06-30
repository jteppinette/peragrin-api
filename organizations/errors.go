package organizations

import (
	"errors"
)

var (
	errOrganizationIDRequired = errors.New("organization id required")
	errCommunityIDRequired    = errors.New("community id required")

	errUploadLogo = errors.New("upload logo")
	errUpdateOrganization = errors.New("update organization")
)
