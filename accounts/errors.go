package accounts

import "errors"

var (
	errAccountIDRequired = errors.New("account id required")

	errAccountNotFound = errors.New("account not found")
)
