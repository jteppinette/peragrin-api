package accounts

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
	"github.com/unrolled/render"
	"gitlab.com/peragrin/api/models"
)

var (
	rend = render.New().JSON
)

// Config defines a single instance of the accounts package.
type Config struct {
	DBClient    *sqlx.DB
	MailClient  *gochimp.MandrillAPI
	Clock       models.Timer
	TokenSecret string
	AppDomain   string
}

// Init generates an accounts.Config instance.
func Init(dbClient *sqlx.DB, mailClient *gochimp.MandrillAPI, clock models.Timer, tokenSecret string, appDomain string) *Config {
	return &Config{dbClient, mailClient, clock, tokenSecret, appDomain}
}
