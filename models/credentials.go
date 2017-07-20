package models

import (
	"database/sql"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/jmoiron/sqlx"
)

// Credentials represents the form necessary to authenticate
// against the Peragrin system.
type Credentials struct {
	Account
	Password string `json:"password"`
}

// Authenticate checks the provided credentials against the database. If
// the provided password passes a hashed comparison against the databases hash,
// then the corresponding account will be returned.
func (c *Credentials) Authenticate(client *sqlx.DB) (*Account, error) {
	credentials, err := GetCredentialsByEmail(c.Email, client)
	if err != nil {
		return nil, err
	}
	if credentials == nil {
		log.WithFields(log.Fields{"email": c.Email}).Info(errAccountNotFound.Error())
		return nil, errAccountNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(c.Password)); err != nil {
		log.WithFields(log.Fields{"email": c.Email, "error": err.Error()}).Info(errInvalidCredentials.Error())
		return nil, errInvalidCredentials
	}
	return &credentials.Account, nil
}

// GetCredentialsByEmail returns the credentials in the database that matches the provided
// email address.
func GetCredentialsByEmail(email string, client *sqlx.DB) (*Credentials, error) {
	c := &Credentials{}
	if err := client.Get(c, "SELECT * FROM Account WHERE LOWER(email) = $1;", strings.ToLower(email)); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return c, nil
}
