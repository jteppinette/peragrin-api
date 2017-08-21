package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
	"golang.org/x/crypto/bcrypt"
)

// Accounts is a slice of account structs.
type Accounts []Account

// Account represents an entity that can login into the Peragrin system.
type Account struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsSuper   bool   `json:"isSuper"`

	// Expiration is used to define the time left for a provided membership.
	// This information is only useful when in the context of a Membership.
	Expiration time.Time `json:"expiration,omitempty"`
}

// Create adds all accounts in the provided slice to the database.
func (accounts *Accounts) Create(client *sqlx.DB) error {
	if len(*accounts) == 0 {
		return nil
	}

	statement := "INSERT INTO Account (email, firstName, lastName) VALUES "
	args := make([]interface{}, len(*accounts)*3)

	for i, account := range *accounts {
		statement = statement + "(?, ?, ?),"
		set := i * 3
		args[set+0] = account.Email
		args[set+1] = account.FirstName
		args[set+2] = account.LastName
	}

	statement = statement[0:len(statement)-1] + " RETURNING id, email, firstName, lastName, isSuper;"
	created := Accounts{}
	if err := client.Select(&created, client.Rebind(statement), args...); err != nil {
		return err
	}
	*accounts = created
	return nil
}

// Save creates or updates the given account in the database.
func (a *Account) Save(client *sqlx.DB) error {
	if a.ID != 0 {
		return client.Get(a, "UPDATE Account SET email = $2, firstName = $3, lastName = $4 WHERE id = $1 RETURNING id, email, firstName, lastName, isSuper;", a.ID, a.Email, a.FirstName, a.LastName)
	}
	return client.Get(a, "INSERT INTO Account (email, firstName, lastName) VALUES ($1, $2, $3) RETURNING id, email, firstName, lastName, isSuper;", a.Email, a.FirstName, a.LastName)
}

// SetPassword sets the account's password.
func (a *Account) SetPassword(password string, client *sqlx.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if _, err := client.Exec("UPDATE Account SET password = $2 WHERE id = $1;", a.ID, string(hash)); err != nil {
		return err
	}
	return nil
}

// SetIsSuper sets the account's super user status.
func (a *Account) SetIsSuper(isSuper bool, client *sqlx.DB) error {
	return client.Get(a, "UPDATE Account SET isSuper = $2 WHERE id = $1 RETURNING id, email, firstName, lastName, isSuper;", a.ID, isSuper)
}

type AuthTokenClaims struct {
	jwt.StandardClaims
	Account
}

func (a Account) AuthToken(key string, expiration time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, AuthTokenClaims{
		jwt.StandardClaims{ExpiresAt: time.Now().Add(expiration).Unix()},
		a,
	}).SignedString([]byte(key))
}

// SendResetPasswordEmail sends a templated reset password email for the provided user.
func (a *Account) SendResetPasswordEmail(appDomain, tokenSecret string, client *gochimp.MandrillAPI) error {
	token, err := a.AuthToken(tokenSecret, time.Hour*24)
	if err != nil {
		return err
	}

	merge := []gochimp.Var{{"RESET_PASSWORD_LINK", fmt.Sprintf("%s/#/auth/set-password?token=%s", appDomain, token)}}
	rendered, err := client.TemplateRender("reset-password", nil, merge)
	if err != nil {
		return err
	}

	if _, err := client.MessageSend(gochimp.Message{
		Html:      rendered,
		Subject:   "Reset Password",
		FromEmail: "donotreply@peragrin.com",
		FromName:  "Peragrin",
		To:        []gochimp.Recipient{{Email: a.Email}},
	}, false); err != nil {
		return err
	}
	return nil
}

func (a *Account) SendActivationEmail(next, appDomain, tokenSecret, name string, client *gochimp.MandrillAPI) error {
	token, err := a.AuthToken(tokenSecret, time.Hour*24*7)
	if err != nil {
		return err
	}

	merge := []gochimp.Var{{"SET_PASSWORD_LINK", fmt.Sprintf("%s/#/auth/activate?token=%s&next=%s", appDomain, token, next)}}
	rendered, err := client.TemplateRender("account-activation", nil, merge)
	if err != nil {
		return err
	}

	var subject string
	if name == "" {
		subject = "Account Activation"
	} else {
		subject = fmt.Sprintf("%s Account Activation", name)
	}

	if _, err := client.MessageSend(gochimp.Message{
		Html:      rendered,
		Subject:   subject,
		FromEmail: "donotreply@peragrin.com",
		FromName:  "Peragrin",
		To:        []gochimp.Recipient{{Email: a.Email}},
	}, false); err != nil {
		return err
	}
	return nil
}

// CreateWithMembership creates a new account with a connection to the
// provided membership.
func (a *Account) CreateWithMembership(membershipID int, client *sqlx.DB) error {
	tx, err := client.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = tx.Get(a, "INSERT INTO Account (email, firstName, lastName) VALUES ($1, $2, $3) RETURNING id, email, firstname, lastName, isSuper;", a.Email, a.FirstName, a.LastName)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO AccountMembership (accountID, membershipID, expiration) VALUES ($1, $2, $3);", a.ID, membershipID, a.Expiration)
	if err != nil {
		return err
	}

	return nil
}

// UpdateWithMembership updates an account and account membership relationship.
func (a *Account) UpdateWithMembership(membershipID int, client *sqlx.DB) error {
	tx, err := client.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = tx.Get(a, "UPDATE Account SET firstName = $2, lastName = $3 WHERE id = $1 RETURNING id, firstname, lastName, isSuper;", a.ID, a.FirstName, a.LastName)
	if err != nil {
		return err
	}

	err = tx.Get(a, "UPDATE AccountMembership SET expiration = $3 WHERE accountID = $1 AND membershipID = $2 RETURNING expiration;", a.ID, membershipID, a.Expiration)
	if err != nil {
		return err
	}

	return nil
}

// CreateWithAccount creates a new account with a connection to the
// provided organization.
func (a *Account) CreateWithOrganization(organizationID int, client *sqlx.DB) error {
	tx, err := client.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = tx.Get(a, "INSERT INTO Account (email, firstName, lastName) VALUES ($1, $2, $3) RETURNING id, email, firstName, lastName, isSuper;", a.Email, a.FirstName, a.LastName)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO AccountOrganization (accountID, organizationID) VALUES ($1, $2);", a.ID, organizationID)
	if err != nil {
		return err
	}

	return nil
}

// AddMembership adds a new membership to the given account.
func (a *Account) AddMembership(membershipID int, client *sqlx.DB) error {
	if _, err := client.Exec("INSERT INTO AccountMembership (accountID, membershipID, expiration) VALUES ($1, $2, $3);", a.ID, membershipID, a.Expiration); err != nil {
		return err
	}
	return nil
}

// RemoveMembership removes a membership from the given account.
func (a *Account) RemoveMembership(membershipID int, client *sqlx.DB) error {
	if _, err := client.Exec("DELETE FROM AccountMembership WHERE accountID = $1 AND membershipID = $2;", a.ID, membershipID); err != nil {
		return err
	}
	return nil
}

// AddOrganization adds a new organization to the given account.
func (a *Account) AddOrganization(organizationID int, client *sqlx.DB) error {
	if _, err := client.Exec("INSERT INTO AccountOrganization (accountID, organizationID) VALUES ($1, $2);", a.ID, organizationID); err != nil {
		return err
	}
	return nil
}

// RemoveOrganization removes an organization from the given account.
func (a *Account) RemoveOrganization(organizationID int, client *sqlx.DB) error {
	if _, err := client.Exec("DELETE FROM AccountOrganization WHERE accountID = $1 AND organizationID = $2;", a.ID, organizationID); err != nil {
		return err
	}
	return nil
}

// GetAccountByEmail returns the account in the database that matches the provided
// email address.
func GetAccountByEmail(email string, client *sqlx.DB) (*Account, error) {
	a := &Account{}
	if err := client.Get(a, "SELECT id, email, firstName, lastName, isSuper FROM Account WHERE LOWER(email) = $1;", strings.ToLower(email)); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return a, nil
}

// GetAccountByID returns the account in the database that matches the provided
// id.
func GetAccountByID(id int, client *sqlx.DB) (*Account, error) {
	a := &Account{}
	if err := client.Get(a, "SELECT id, email, firstName, lastName, isSuper FROM Account WHERE id = $1;", id); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return a, nil
}

// GetAccountsByMembership returns all accounts with the provided membership.
func GetAccountsByMembership(membershipID int, client *sqlx.DB) (Accounts, error) {
	accounts := Accounts{}
	if err := client.Select(&accounts, `
		SELECT Account.id, Account.email, Account.firstName, Account.lastName, Account.isSuper, AccountMembership.Expiration
		FROM Account INNER JOIN AccountMembership ON (Account.id = AccountMembership.accountID)
		WHERE AccountMembership.membershipID = $1
	`, membershipID); err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountsByOrganization returns all accounts that are operating the provided organization.
func GetAccountsByOrganization(organizationID int, client *sqlx.DB) (Accounts, error) {
	accounts := Accounts{}
	if err := client.Select(&accounts, `
		SELECT Account.id, Account.email, Account.firstName, Account.lastName, Account.isSuper
		FROM Account INNER JOIN AccountOrganization ON (Account.id = AccountOrganization.accountID)
		WHERE AccountOrganization.organizationID = $1
	`, organizationID); err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountsByEmails returns the set of accounts that have an email in the provided list.
func GetAccountsByEmails(emails []string, client *sqlx.DB) (Accounts, error) {
	accounts := Accounts{}
	query, args, err := sqlx.In(`
		SELECT Account.id, Account.email, Account.firstName, Account.lastName, Account.isSuper
		FROM Account WHERE email IN (?);
	`, emails)
	if err != nil {
		return nil, err
	}
	if err := client.Select(&accounts, client.Rebind(query), args...); err != nil {
		return nil, err
	}
	return accounts, nil
}
