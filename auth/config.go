package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattbaird/gochimp"
	minio "github.com/minio/minio-go"
	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

type Config struct {
	DBClient         *sqlx.DB
	StoreClient      *minio.Client
	TokenSecret      string
	Clock            timer
	LocationIQAPIKey string
	AppDomain        string
	MailClient       *gochimp.MandrillAPI
}

func Init(dbClient *sqlx.DB, storeClient *minio.Client, tokenSecret string, locationIQAPIKey string, appDomain string, mailClient *gochimp.MandrillAPI) *Config {
	return &Config{dbClient, storeClient, tokenSecret, clock{}, locationIQAPIKey, appDomain, mailClient}
}
