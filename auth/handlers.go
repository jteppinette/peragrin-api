package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/pkg/errors"
	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// GetAccountHandler returns the currently authenticated account. To function properly,
// a preceding middleware must add the "account" key to the request context.
func (c *Config) GetAccountHandler(r *http.Request) *service.Response {
	if account, ok := context.GetOk(r, "account"); ok {
		return service.NewResponse(nil, http.StatusOK, account)
	}
	return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
}

// ListOrganizationsHandler generates a response object containing the organizations that are
// operated by the currently authenticated account.
func (c *Config) ListOrganizationsHandler(r *http.Request) *service.Response {
	account, ok := context.Get(r, "account").(models.Account)
	if !ok {
		return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
	}

	organizations, err := models.GetOrganizationsByAccount(account.ID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	// TODO: Mock out the static store.
	if c.StoreClient != nil {
		if err := organizations.SetPresignedLogoLinks(c.StoreClient); err != nil {
			return service.NewResponse(err, http.StatusBadRequest, nil)
		}
	}

	return service.NewResponse(nil, http.StatusOK, organizations)
}

// CreateOrganizationHandler saves a new organization to the database.
func (c *Config) CreateOrganizationHandler(r *http.Request) *service.Response {
	account, ok := context.Get(r, "account").(models.Account)
	if !ok {
		return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
	}

	organization := models.Organization{}
	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	// If there is a geocode lookup failure, then log the failure. We
	// will just let the user manually enter the coordinates.
	if err := organization.SetGeo(c.LocationIQAPIKey); err != nil {
		logrus.WithFields(logrus.Fields{
			"street":  organization.Street,
			"city":    organization.City,
			"state":   organization.State,
			"country": organization.Country,
			"zip":     organization.Zip,
		}).Error(errors.Wrap(err, errGeocode.Error()))
	}

	if err := organization.Create(account.ID, c.DBClient); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreateOrganization.Error()), http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusCreated, organization)
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

	str, err := token(c.TokenSecret, account, c.Clock)
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
	if err := a.Save(c.DBClient); err != nil {
		return service.NewResponse(errors.Wrap(err, errRegistrationFailed.Error()), http.StatusBadRequest, nil)
	}

	str, err := token(c.TokenSecret, a, c.Clock)
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
