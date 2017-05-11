package posts

import (
	"errors"
)

var (
	errCreatePost             = errors.New("create post")
	errAuthenticationRequired = errors.New("authentication required")
	errOrganizationIDRequired = errors.New("organization id required")
)
