package auth

import (
	log "github.com/Sirupsen/logrus"
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
func (creds Credentials) Authenticate(c *Config, requestID string) (models.Account, error) {
	account, err := models.GetAccountByEmail(creds.Email, c.DBClient)
	if err != nil {
		log.WithFields(log.Fields{
			"email": creds.Email, "error": err.Error(), "id": requestID,
		}).Info(errAccountNotFound.Error())
		return models.Account{}, errAccountNotFound
	}
	if err := account.ValidatePassword(creds.Password); err != nil {
		log.WithFields(log.Fields{
			"email": creds.Email, "error": err.Error(), "id": requestID,
		}).Info(errInvalidCredentials.Error())
		return models.Account{}, errInvalidCredentials
	}
	// Do not allow the hashed password to be returned outside of this function.
	account.Password = ""
	return *account, nil
}
