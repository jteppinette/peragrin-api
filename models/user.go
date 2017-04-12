package models

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Users []User

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) Save(client *sqlx.DB) error {
	if u.ID != "" {
		return client.Get(u, "UPDATE users SET username = $2, password = $3 WHERE id = $1 RETURNING *;", u.ID, u.Username, u.Password)
	} else {
		return client.Get(u, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *;", u.Username, u.Password)
	}
}

func ListUsers(client *sqlx.DB) (Users, error) {
	users := Users{}
	if err := client.Select(&users, "SELECT * FROM users;"); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByUsername(username string, client *sqlx.DB) (*User, error) {
	u := &User{}
	if err := client.Get(u, "SELECT * FROM users WHERE username = $1;", username); err != nil {
		return nil, err
	}
	return u, nil
}
