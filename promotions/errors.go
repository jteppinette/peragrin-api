package promotions

import (
	"errors"
)

var (
	errPromotionIDRequired    = errors.New("promotion id required")
	errAuthenticationRequired = errors.New("authentication required")
)
