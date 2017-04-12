package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

func (c *Config) UserHandler(w http.ResponseWriter, r *http.Request) {
	if user, ok := context.GetOk(r, "user"); ok {
		rend(w, http.StatusOK, user)
		return
	}
	service.Error(w, http.StatusUnauthorized, errAuthenticationRequired)
}

func (c *Config) LoginHandler(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		service.Error(w, http.StatusBadRequest, errBadCredentialsFormat)
		return
	}

	user, err := creds.Authenticate(c)
	if err != nil {
		service.Error(w, http.StatusUnauthorized, err)
		return
	}

	str, err := token(c.TokenSecret, user)
	if err != nil {
		service.Error(w, http.StatusUnauthorized, err)
		return
	}

	rend(w, http.StatusOK, struct {
		Token string `json:"token"`
		models.User
	}{str, user})
}

func (c *Config) RequiredMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			service.Error(w, http.StatusUnauthorized, errAuthenticationRequired)
			return
		}

		var user models.User
		if strings.HasPrefix(authorization, "Bearer ") {
			token, err := jwt.ParseWithClaims(strings.Split(authorization, " ")[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(c.TokenSecret), nil
			})
			if err != nil {
				service.Error(w, http.StatusUnauthorized, fmt.Errorf("%v: %v", errJWTAuth, err))
				return
			}
			user = token.Claims.(*Claims).User
		} else {
			username, password, ok := r.BasicAuth()
			if !ok {
				service.Error(w, http.StatusBadRequest, fmt.Errorf("%v: %v", errBasicAuth, errBadCredentialsFormat))
				return
			}
			var err error
			user, err = Credentials{username, password}.Authenticate(c)
			if err != nil {
				service.Error(w, http.StatusUnauthorized, fmt.Errorf("%v: %v", errBasicAuth, err))
				return
			}
		}

		context.Set(r, "user", user)
		h.ServeHTTP(w, r)
	})
}
