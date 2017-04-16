package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"gitlab.com/peragrin/api/models"
)

// UserHandler returns the currently authenticated user. To function properly,
// a preceding middleware must add the "user" key to the request context.
func (c *Config) UserHandler(w http.ResponseWriter, r *http.Request) {
	if user, ok := context.GetOk(r, "user"); ok {
		rend(w, http.StatusOK, user)
		return
	}
	// errAuthenticateRequired
	rend(w, http.StatusUnauthorized, nil)
}

type authUser struct {
	Token string `json:"token"`
	models.User
}

// LoginHandler reads JSON encoded username and password from the provided request
// and attempts to authenticate these credentials.
// If succesful, an authUser object will be returned to the client.
func (c *Config) LoginHandler(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// errBadCredentialsFormat
		rend(w, http.StatusBadRequest, nil)
		return
	}

	user, err := creds.Authenticate(c)
	if err != nil {
		// err
		rend(w, http.StatusUnauthorized, nil)
		return
	}

	str, err := token(c.TokenSecret, user)
	if err != nil {
		// err
		rend(w, http.StatusUnauthorized, nil)
		return
	}

	rend(w, http.StatusOK, authUser{str, user})
}

// RequireAuthMiddleware attempts to authenticate the incoming request using
// Basic and JWT authentication strategies. If successful, a "user" key will be
// added to the request context. Otherwise, an HTTP Unauthorized will be
// returned to the client.
func (c *Config) RequireAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			// errAuthenticationRequired
			rend(w, http.StatusUnauthorized, nil)
			return
		}

		var user models.User
		if strings.HasPrefix(authorization, "Bearer ") {
			token, err := jwt.ParseWithClaims(strings.Split(authorization, " ")[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(c.TokenSecret), nil
			})
			if err != nil {
				// fmt.Errorf("%v: %v", errJWTAuth, err)
				rend(w, http.StatusUnauthorized, nil)
				return
			}
			user = token.Claims.(*Claims).User
		} else {
			username, password, ok := r.BasicAuth()
			if !ok {
				// fmt.Errorf("%v: %v", errBasicAuth, errBadCredentialsFormat))
				rend(w, http.StatusBadRequest, nil)
				return
			}
			var err error
			user, err = Credentials{username, password}.Authenticate(c)
			if err != nil {
				// fmt.Errorf("%v: %v", errBasicAuth, err)
				rend(w, http.StatusUnauthorized, nil)
				return
			}
		}

		context.Set(r, "user", user)
		h.ServeHTTP(w, r)
	})
}
