package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"

	"github.com/jteppinette/peragrin-api/mail"
)

var (
	rend = render.New().JSON
)

// Config defines a single instance of the auth package.
type Config struct {
	DBClient    *sqlx.DB
	MailClient  *mail.Config
	TokenSecret string
	AppDomain   string
}

// Init generates an auth.Config instance.
func Init(dbClient *sqlx.DB, mailClient *mail.Config, tokenSecret, appDomain string) *Config {
	return &Config{dbClient, mailClient, tokenSecret, appDomain}
}
