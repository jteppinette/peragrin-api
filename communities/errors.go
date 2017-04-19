package communities

import (
	"errors"
)

var (
	errListCommunities   = errors.New("list communities")
	errListOrganizations = errors.New("list organizations")

	errCommunityIDRequired = errors.New("community id required")
)
