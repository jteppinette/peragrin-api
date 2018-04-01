package organizations

import (
	"github.com/jmoiron/sqlx"
	minio "github.com/minio/minio-go"

	"github.com/jteppinette/peragrin-api/mail"
)

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	DBClient    *sqlx.DB
	StoreClient *minio.Client
	MailClient  *mail.Config
	TokenSecret string
	AppDomain   string
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(dbClient *sqlx.DB, storeClient *minio.Client, mailClient *mail.Config, tokenSecret, appDomain string) *Config {
	return &Config{dbClient, storeClient, mailClient, tokenSecret, appDomain}
}
