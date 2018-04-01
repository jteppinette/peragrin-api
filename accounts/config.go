package accounts

import (
	"github.com/jmoiron/sqlx"
	minio "github.com/minio/minio-go"
	"github.com/unrolled/render"

	"github.com/jteppinette/peragrin-api/mail"
)

var (
	rend = render.New().JSON
)

// Config defines a single instance of the accounts package.
type Config struct {
	DBClient    *sqlx.DB
	StoreClient *minio.Client
	MailClient  *mail.Config
	TokenSecret string
	AppDomain   string
}

// Init generates an accounts.Config instance.
func Init(dbClient *sqlx.DB, storeClient *minio.Client, mailClient *mail.Config, tokenSecret, appDomain string) *Config {
	return &Config{dbClient, storeClient, mailClient, tokenSecret, appDomain}
}
