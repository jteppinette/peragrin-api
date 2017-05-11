package auth

import (
	"errors"
)

var (
	errBadCredentialsFormat               = errors.New("bad credentials format")
	errInvalidCredentials                 = errors.New("invalid credentials")
	errAccountNotFound                    = errors.New("account not found")
	errAuthenticationRequired             = errors.New("authentication required")
	errAuthenticationStrategyNotSupported = errors.New("authentication strategy not supported")
	errRegistrationFailed                 = errors.New("registration failed")
	errGeocodeFailed                      = errors.New("geocode failed")

	errBasicAuth = errors.New("basic auth")
	errJWTAuth   = errors.New("jwt auth")
)
