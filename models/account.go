package models

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Accounts is a slice of account structs.
type Accounts []Account

// Account represents an entity that can login into the Peragrin system.
type Account struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// ValidatePassword compares the given password to the password hash
// that is in the receiving account struct.
func (a *Account) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
}

// SetPassword generates a password hash using the bcrypt algorithm.
// This hash is then stored on the receiving account struct in the password field.
func (a *Account) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hash)
	return nil
}

// Save creates or updates the given account in the database.
func (a *Account) Save(client *sqlx.DB) error {
	if a.ID != 0 {
		return client.Get(a, "UPDATE Account SET email = $2, password = $3 WHERE id = $1 RETURNING *;", a.ID, a.Email, a.Password)
	}
	return client.Get(a, "INSERT INTO Account (email, password) VALUES ($1, $2) RETURNING *;", a.Email, a.Password)
}

// AddOperator creates a relationship that defines this account as an
// an operator of the provided organization.
func (a *Account) AddOperator(organizationID int, client *sqlx.DB) error {
	o := Operator{AccountID: a.ID, OrganizationID: organizationID}
	return o.Save(client)
}

// ListAccounts returns all accounts in the database.
func ListAccounts(client *sqlx.DB) (Accounts, error) {
	accounts := Accounts{}
	if err := client.Select(&accounts, "SELECT * FROM Account;"); err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountByEmail returns the account in the database that matches the provided
// email address. If there is not matching account, then an error is returned.
func GetAccountByEmail(email string, client *sqlx.DB) (*Account, error) {
	a := &Account{}
	if err := client.Get(a, "SELECT * FROM Account WHERE email = $1;", email); err != nil {
		return nil, err
	}
	return a, nil
}
