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

// Authenticate retrieves an account object using the provided credentials.
// If the password hashes validate, then the account will be returned.
func (creds Credentials) Authenticate(c *Config) (models.Account, error) {
	account, err := models.GetAccountByEmail(creds.Email, c.DBClient)
	if err != nil {
		return models.Account{}, errors.Wrap(err, errAccountNotFound.Error())
	}
	if err := account.ValidatePassword(creds.Password); err != nil {
		return models.Account{}, errors.Wrap(err, errInvalidCredentials.Error())
	}
	// Do not allow the hashed password to be returned outside of this function.
	account.Password = ""
	return *account, nil
}
