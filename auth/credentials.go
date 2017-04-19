package auth

import (
	"github.com/pkg/errors"
	"gitlab.com/peragrin/api/models"
)

// Credentials represents the form necessary to authenticate
// against the Peragrin system.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Authenticate retrieves a user object using the provided credentials.
// If the password hashes validate, then the user will be returned.
func (creds Credentials) Authenticate(c *Config) (models.User, error) {
	user, err := models.GetUserByEmail(creds.Email, c.Client)
	if err != nil {
		return models.User{}, errors.Wrap(err, errUserNotFound.Error())
	}
	if err := user.ValidatePassword(creds.Password); err != nil {
		return models.User{}, errors.Wrap(err, errInvalidCredentials.Error())
	}
	// Do not allow the hashed password to be returned outside of this function.
	user.Password = ""
	return *user, nil
}
