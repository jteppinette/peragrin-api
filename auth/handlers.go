package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/pkg/errors"
	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// SetPasswordHandler allows an authenticated user to set a new password.
func (c *Config) SetPasswordHandler(r *http.Request) *service.Response {
	account, ok := context.Get(r, "account").(models.Account)
	if !ok {
		return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
	}
	form := struct {
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := account.SetPassword(form.Password, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}

// ActivateHandler allows an authenticated user to set a new password.
func (c *Config) ActivateHandler(r *http.Request) *service.Response {
	account, ok := context.Get(r, "account").(models.Account)
	if !ok {
		return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
	}
	form := struct {
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := account.SetPassword(form.Password, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	str, err := account.AuthToken(c.TokenSecret, time.Hour*24)
	if err != nil {
		return service.NewResponse(err, http.StatusUnauthorized, map[string]string{"msg": err.Error()})
	}

	return service.NewResponse(nil, http.StatusOK, struct {
		Token string `json:"token"`
	}{str})
}

// ForgotPasswordHandler generates a token that can be used to reset the password
// of the account with the provided email address.
func (c *Config) ForgotPasswordHandler(r *http.Request) *service.Response {
	form := struct {
		Email string `json:"email"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	account, err := models.GetAccountByEmail(form.Email, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	if account == nil {
		// We don't want to allow account enumeration, so we are just goign to
		// log the error and return the success response.
		log.WithFields(log.Fields{
			"email": form.Email,
		}).Info(errAccountNotFound.Error())
		return service.NewResponse(nil, http.StatusOK, nil)
	}

	if err := account.SendResetPasswordEmail(c.AppDomain, c.TokenSecret, c.MailClient); err != nil {
		return service.NewResponse(err, http.StatusInternalServerError, nil)
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}

// LoginHandler reads a JSON encoded email and password from the provided request
// and attempts to authenticate these credentials.
// If succesful, a token object will be returned to the client.
func (c *Config) LoginHandler(r *http.Request) *service.Response {
	creds := models.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		return service.NewResponse(errors.Wrap(err, errBadCredentialsFormat.Error()), http.StatusBadRequest, nil)
	}

	account, err := creds.Authenticate(c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusUnauthorized, nil)
	}

	str, err := account.AuthToken(c.TokenSecret, time.Hour*24)
	if err != nil {
		return service.NewResponse(err, http.StatusUnauthorized, map[string]string{"msg": err.Error()})
	}

	return service.NewResponse(nil, http.StatusOK, struct {
		Token string `json:"token"`
	}{str})
}

// RegisterHandler creates a new account and returns a account object.
func (c *Config) RegisterHandler(r *http.Request) *service.Response {
	account := models.Account{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return service.NewResponse(nil, http.StatusBadRequest, nil)
	}

	// If the account already existing then simply return an HTTP 200.
	if existing, err := models.GetAccountByEmail(account.Email, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	} else if existing != nil {
		log.WithFields(log.Fields{
			"email":     existing.Email,
			"accountID": existing.ID,
		}).Info("registration: account already exists")
		return service.NewResponse(nil, http.StatusOK, nil)
	}

	// If an error occurs while creating this new account, return a standard HTTP 200.
	// This is to prevent any account enumeration.
	if err := account.Save(c.DBClient); err != nil {
		log.WithFields(log.Fields{
			"email": account.Email, "error": err.Error(), "id": r.Header.Get("X-Request-ID"),
		}).Info(errRegistration.Error())
		return service.NewResponse(nil, http.StatusOK, nil)
	}

	if err := account.SendActivationEmail("", c.AppDomain, c.TokenSecret, "", c.MailClient); err != nil {
		log.WithFields(log.Fields{
			"email": account.Email, "error": err.Error(), "id": r.Header.Get("X-Request-ID"),
		}).Info(errAccountActivationEmail.Error())
	}

	return service.NewResponse(nil, http.StatusOK, nil)
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

		var account *models.Account
		if strings.HasPrefix(authorization, "Bearer ") {
			token, err := jwt.ParseWithClaims(strings.Split(authorization, " ")[1], &models.AuthTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(c.TokenSecret), nil
			})
			if err != nil {
				return service.NewResponse(errors.Wrap(err, errJWTAuth.Error()), http.StatusUnauthorized, nil)
			}
			account = &token.Claims.(*models.AuthTokenClaims).Account
		} else if strings.HasPrefix(authorization, "Basic ") {
			email, password, ok := r.BasicAuth()
			if !ok {
				return service.NewResponse(errors.Wrap(errBadCredentialsFormat, errBasicAuth.Error()), http.StatusBadRequest, nil)
			}
			var err error
			credentials := &models.Credentials{models.Account{Email: email}, password}
			account, err = credentials.Authenticate(c.DBClient)
			if err != nil {
				return service.NewResponse(errors.Wrap(err, errBasicAuth.Error()), http.StatusUnauthorized, nil)
			}
		} else {
			return service.NewResponse(errAuthenticationStrategyNotSupported, http.StatusUnauthorized, nil)
		}

		context.Set(r, "account", *account)
		return h(r)
	}
}
