package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

func (c *Config) UserHandler(w http.ResponseWriter, r *http.Request) {
	if user, ok := context.GetOk(r, "user"); ok {
		json(w, http.StatusOK, user)
		return
	}
	service.Error(w, http.StatusUnauthorized, errAuthenticationRequired)
}

func (c *Config) RequiredMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			service.Error(w, http.StatusBadRequest, errValidBasicAuthCredentialsRequired)
			return
		}
		user, err := models.GetUserByUsername(username, c.Client)
		if err != nil {
			service.Error(w, http.StatusBadRequest, fmt.Errorf("%+v: %+v", errUserNotFound, err))
			return
		}
		if err := user.ValidatePassword(password); err != nil {
			service.Error(w, http.StatusUnauthorized, errInvalidCredentials)
			return
		}
		context.Set(r, "user", user)

		h.ServeHTTP(w, r)
	})
}
