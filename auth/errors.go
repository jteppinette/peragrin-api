package auth

import (
	"errors"
)

var (
	errBadCredentialsFormat   = errors.New("bad credentials format")
	errInvalidCredentials     = errors.New("invalid credentials")
	errUserNotFound           = errors.New("user not found")
	errAuthenticationRequired = errors.New("authentication required")

	errBasicAuth = errors.New("basic auth")
	errJWTAuth   = errors.New("jwt auth")
)
