package communities

import (
	"errors"
)

var (
	errCommunityIDRequired = errors.New("community id required")
	errCreateOrganization  = errors.New("create organization")
	errGeocodeFailed       = errors.New("geocode failed")
)
