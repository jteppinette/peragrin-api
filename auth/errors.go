package auth

import (
	"errors"
)

var (
	errInvalidCredentials                = errors.New("invalid credentials")
	errValidBasicAuthCredentialsRequired = errors.New("valid basic auth credentials required")
	errUserNotFound                      = errors.New("user not found")
	errAuthenticationRequired            = errors.New("authentication required")
)
