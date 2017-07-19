package accounts

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
	minio "github.com/minio/minio-go"
	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

// Config defines a single instance of the accounts package.
type Config struct {
	DBClient         *sqlx.DB
	StoreClient      *minio.Client
	MailClient       *gochimp.MandrillAPI
	TokenSecret      string
	AppDomain        string
	LocationIQAPIKey string
}

// Init generates an accounts.Config instance.
func Init(dbClient *sqlx.DB, storeClient *minio.Client, mailClient *gochimp.MandrillAPI, tokenSecret, appDomain, locationIQAPIKey string) *Config {
	return &Config{dbClient, storeClient, mailClient, tokenSecret, appDomain, locationIQAPIKey}
}
