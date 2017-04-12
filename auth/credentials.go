package auth

import (
	"fmt"

	"gitlab.com/peragrin/api/models"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (creds Credentials) Authenticate(c *Config) (models.User, error) {
	user, err := models.GetUserByUsername(creds.Username, c.Client)
	if err != nil {
		return models.User{}, fmt.Errorf("%+v: %+v", errUserNotFound, err)
	}
	if err := user.ValidatePassword(creds.Password); err != nil {
		return models.User{}, fmt.Errorf("%+v: %+v", errInvalidCredentials, err)
	}
	return *user, nil
}
