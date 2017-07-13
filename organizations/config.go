package organizations

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
	minio "github.com/minio/minio-go"
	"gitlab.com/peragrin/api/models"
)

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	DBClient    *sqlx.DB
	StoreClient *minio.Client
	MailClient  *gochimp.MandrillAPI
	Clock       models.Timer
	TokenSecret string
	AppDomain   string
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(dbClient *sqlx.DB, storeClient *minio.Client, mailClient *gochimp.MandrillAPI, clock models.Timer, tokenSecret, appDomain string) *Config {
	return &Config{dbClient, storeClient, mailClient, clock, tokenSecret, appDomain}
}
