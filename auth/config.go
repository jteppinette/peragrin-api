package auth

import (
	"github.com/jmoiron/sqlx"
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
	MandrillKey      string
}

func Init(dbClient *sqlx.DB, storeClient *minio.Client, tokenSecret, locationIQAPIKey, appDomain, mandrillKey string) *Config {
	return &Config{dbClient, storeClient, tokenSecret, clock{}, locationIQAPIKey, appDomain, mandrillKey}
}
