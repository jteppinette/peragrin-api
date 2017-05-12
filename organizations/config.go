package organizations

import (
	"github.com/jmoiron/sqlx"
)

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	Client           *sqlx.DB
	LocationIQAPIKey string
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(client *sqlx.DB, locationIQAPIKey string) *Config {
	return &Config{client, locationIQAPIKey}
}
