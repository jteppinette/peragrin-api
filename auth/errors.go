package auth

import "errors"

var (
	errBadCredentialsFormat               = errors.New("bad credentials format")
	errAccountNotFound                    = errors.New("account not found")
	errAuthenticationRequired             = errors.New("authentication required")
	errAuthenticationStrategyNotSupported = errors.New("authentication strategy not supported")
	errRegistration                       = errors.New("registration")
	errAccountActivationEmail             = errors.New("account activation email")

	errBasicAuth = errors.New("basic auth")
	errJWTAuth   = errors.New("jwt auth")
)
