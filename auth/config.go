package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
	minio "github.com/minio/minio-go"
	"github.com/unrolled/render"
	"gitlab.com/peragrin/api/models"
)

var (
	rend = render.New().JSON
)

// Config defines a single instance of the auth package.
type Config struct {
	DBClient         *sqlx.DB
	StoreClient      *minio.Client
	MailClient       *gochimp.MandrillAPI
	Clock            models.Timer
	TokenSecret      string
	LocationIQAPIKey string
	AppDomain        string
}

// Init generates an auth.Config instance.
func Init(dbClient *sqlx.DB, storeClient *minio.Client, mailClient *gochimp.MandrillAPI, clock models.Timer, tokenSecret string, locationIQAPIKey string, appDomain string) *Config {
	return &Config{dbClient, storeClient, mailClient, clock, tokenSecret, locationIQAPIKey, appDomain}
}
