package memberships

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
)

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	DBClient    *sqlx.DB
	MailClient  *gochimp.MandrillAPI
	TokenSecret string
	AppDomain   string
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(dbClient *sqlx.DB, mailClient *gochimp.MandrillAPI, tokenSecret, appDomain string) *Config {
	return &Config{dbClient, mailClient, tokenSecret, appDomain}
}
