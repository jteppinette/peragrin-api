package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/peragrin/api/models"
)

type Claims struct {
	jwt.StandardClaims
	models.User
}

func token(key string, user models.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24).Unix()},
		user,
	}).SignedString([]byte(key))
}
