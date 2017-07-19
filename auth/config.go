package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

// Config defines a single instance of the auth package.
type Config struct {
	DBClient    *sqlx.DB
	MailClient  *gochimp.MandrillAPI
	TokenSecret string
	AppDomain   string
}

// Init generates an auth.Config instance.
func Init(dbClient *sqlx.DB, mailClient *gochimp.MandrillAPI, tokenSecret, appDomain string) *Config {
	return &Config{dbClient, mailClient, tokenSecret, appDomain}
}
