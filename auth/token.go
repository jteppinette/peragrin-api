package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/peragrin/api/models"
)

type timer interface {
	Now() time.Time
}

type clock struct{}

func (clock) Now() time.Time {
	return time.Now()
}

// Claims is the struct that is stored inside a JWT.
type Claims struct {
	jwt.StandardClaims
	models.User
}

func token(key string, user models.User, c timer) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		jwt.StandardClaims{ExpiresAt: c.Now().Add(time.Hour * 24).Unix()},
		user,
	}).SignedString([]byte(key))
}
