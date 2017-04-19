package models

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Users is a slice of user structs.
type Users []User

// User represents an entity that can login into the Peragrin system
// and manage a organization or organizations.
type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	Password       string `json:"-"`
	OrganizationID int    `json:"organizationID"`
}

// ValidatePassword compares the given password to the password hash
// that is in the receiving user struct.
func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// SetPassword generates a password hash using the bcrypt algorithm.
// This hash is then stored on the receiving user struct in the password field.
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// Save creates or updates the given user in the database.
func (u *User) Save(client *sqlx.DB) error {
	if u.ID != 0 {
		return client.Get(u, "UPDATE users SET email = $2, password = $3, organizationID = $4 WHERE id = $1 RETURNING *;", u.ID, u.Email, u.Password, u.OrganizationID)
	} else {
		return client.Get(u, "INSERT INTO users (email, password, organizationID) VALUES ($1, $2, $3) RETURNING *;", u.Email, u.Password, u.OrganizationID)
	}
}

// ListUsers returns all users in the database.
func ListUsers(client *sqlx.DB) (Users, error) {
	users := Users{}
	if err := client.Select(&users, "SELECT * FROM users;"); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByEmail returns the user in the database that matches the provided
// email address. If there is not matching user, then an error is returned.
func GetUserByEmail(email string, client *sqlx.DB) (*User, error) {
	u := &User{}
	if err := client.Get(u, "SELECT * FROM users WHERE email = $1;", email); err != nil {
		return nil, err
	}
	return u, nil
}
