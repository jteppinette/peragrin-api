package organizations

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
	minio "github.com/minio/minio-go"
)

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	DBClient         *sqlx.DB
	StoreClient      *minio.Client
	MailClient       *gochimp.MandrillAPI
	TokenSecret      string
	AppDomain        string
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(dbClient *sqlx.DB, storeClient *minio.Client, mailClient *gochimp.MandrillAPI, tokenSecret, appDomain string) *Config {
	return &Config{dbClient, storeClient, mailClient, tokenSecret, appDomain}
}
