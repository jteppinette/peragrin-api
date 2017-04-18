package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/pkg/errors"
	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// UserHandler returns the currently authenticated user. To function properly,
// a preceding middleware must add the "user" key to the request context.
func (c *Config) UserHandler(r *http.Request) *service.Response {
	if user, ok := context.GetOk(r, "user"); ok {
		return service.NewResponse(nil, http.StatusOK, user)
	}
	return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
}

type authUser struct {
	Token string `json:"token"`
	models.User
}

// LoginHandler reads JSON encoded username and password from the provided request
// and attempts to authenticate these credentials.
// If succesful, an authUser object will be returned to the client.
func (c *Config) LoginHandler(r *http.Request) *service.Response {
	creds := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errBadCredentialsFormat.Error()), http.StatusBadRequest, nil)
	}

	user, err := creds.Authenticate(c)
	if err != nil {
		return service.NewResponse(err, http.StatusUnauthorized, nil)
	}

	str, err := token(c.TokenSecret, user, c.Clock)
	if err != nil {
		return service.NewResponse(err, http.StatusUnauthorized, nil)
	}

	return service.NewResponse(nil, http.StatusOK, authUser{str, user})
}

// RequiredMiddleware attempts to authenticate the incoming request using
// Basic and JWT authentication strategies. If successful, a "user" key will be
// added to the request context. Otherwise, an HTTP Unauthorized will be
// returned to the client.
func (c *Config) RequiredMiddleware(h service.Handler) service.Handler {
	return func(r *http.Request) *service.Response {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
		}

		var user models.User
		if strings.HasPrefix(authorization, "Bearer ") {
			token, err := jwt.ParseWithClaims(strings.Split(authorization, " ")[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(c.TokenSecret), nil
			})
			if err != nil {
				return service.NewResponse(errors.Wrap(err, errJWTAuth.Error()), http.StatusUnauthorized, nil)
			}
			user = token.Claims.(*Claims).User
		} else if strings.HasPrefix(authorization, "Basic ") {
			username, password, ok := r.BasicAuth()
			if !ok {
				return service.NewResponse(errors.Wrap(errBadCredentialsFormat, errBasicAuth.Error()), http.StatusBadRequest, nil)
			}
			var err error
			user, err = Credentials{username, password}.Authenticate(c)
			if err != nil {
				return service.NewResponse(errors.Wrap(err, errBasicAuth.Error()), http.StatusUnauthorized, nil)
			}
		} else {
			return service.NewResponse(errAuthenticationStrategyNotSupported, http.StatusUnauthorized, nil)
		}

		context.Set(r, "user", user)
		return h(r)
	}
}
