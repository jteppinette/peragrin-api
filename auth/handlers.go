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

// AccountHandler returns the currently authenticated account. To function properly,
// a preceding middleware must add the "account" key to the request context.
func (c *Config) AccountHandler(r *http.Request) *service.Response {
	if account, ok := context.GetOk(r, "account"); ok {
		return service.NewResponse(nil, http.StatusOK, account)
	}
	return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
}

// LoginHandler reads a JSON encoded email and password from the provided request
// and attempts to authenticate these credentials.
// If succesful, a token object will be returned to the client.
func (c *Config) LoginHandler(r *http.Request) *service.Response {
	creds := Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		return service.NewResponse(errors.Wrap(err, errBadCredentialsFormat.Error()), http.StatusBadRequest, nil)
	}

	account, err := creds.Authenticate(c)
	if err != nil {
		return service.NewResponse(err, http.StatusUnauthorized, nil)
	}

	str, err := token(c.TokenSecret, account, c.MapboxAPIKey, c.Clock)
	if err != nil {
		return service.NewResponse(err, http.StatusUnauthorized, nil)
	}

	return service.NewResponse(nil, http.StatusOK, struct {
		Token string `json:"token"`
	}{str})
}

// RegisterHandler creates a new account and returns a account object.
func (c *Config) RegisterHandler(r *http.Request) *service.Response {
	creds := Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		return service.NewResponse(errors.Wrap(err, errBadCredentialsFormat.Error()), http.StatusBadRequest, nil)
	}

	a := models.Account{Email: creds.Email}
	a.SetPassword(creds.Password)
	if err := a.Save(c.Client); err != nil {
		return service.NewResponse(errors.Wrap(err, errRegistrationFailed.Error()), http.StatusBadRequest, nil)
	}

	str, err := token(c.TokenSecret, a, c.MapboxAPIKey, c.Clock)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, struct {
		Token string `json:"token"`
	}{str})
}

// RequiredMiddleware attempts to authenticate the incoming request using
// Basic and JWT authentication strategies. If successful, an "account" key will be
// added to the request context. Otherwise, an HTTP Unauthorized will be
// returned to the client.
func (c *Config) RequiredMiddleware(h service.Handler) service.Handler {
	return func(r *http.Request) *service.Response {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
		}

		var account models.Account
		if strings.HasPrefix(authorization, "Bearer ") {
			token, err := jwt.ParseWithClaims(strings.Split(authorization, " ")[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(c.TokenSecret), nil
			})
			if err != nil {
				return service.NewResponse(errors.Wrap(err, errJWTAuth.Error()), http.StatusUnauthorized, nil)
			}
			account = token.Claims.(*Claims).Account
		} else if strings.HasPrefix(authorization, "Basic ") {
			email, password, ok := r.BasicAuth()
			if !ok {
				return service.NewResponse(errors.Wrap(errBadCredentialsFormat, errBasicAuth.Error()), http.StatusBadRequest, nil)
			}
			var err error
			account, err = Credentials{email, password}.Authenticate(c)
			if err != nil {
				return service.NewResponse(errors.Wrap(err, errBasicAuth.Error()), http.StatusUnauthorized, nil)
			}
		} else {
			return service.NewResponse(errAuthenticationStrategyNotSupported, http.StatusUnauthorized, nil)
		}

		context.Set(r, "account", account)
		return h(r)
	}
}
