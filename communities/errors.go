package communities

import (
	"errors"
)

var (
	errListCommunities   = errors.New("list communities")
	errListOrganizations = errors.New("list organizations")
	errListPosts         = errors.New("list posts")

	errCommunityIDRequired = errors.New("community id required")
)
