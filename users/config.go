package users

import (
	"github.com/jmoiron/sqlx"
)

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	Client *sqlx.DB
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(client *sqlx.DB) *Config {
	return &Config{client}
}
