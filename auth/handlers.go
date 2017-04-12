package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
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

	rend(w, http.StatusOK, user)
}

func (c *Config) RequiredMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			service.Error(w, http.StatusBadRequest, errBadCredentialsFormat)
			return
		}

		user, err := Credentials{username, password}.Authenticate(c)
		if err != nil {
			service.Error(w, http.StatusUnauthorized, err)
		}

		context.Set(r, "user", user)
		h.ServeHTTP(w, r)
	})
}
