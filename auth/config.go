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
}

func Init(dbClient *sqlx.DB, storeClient *minio.Client, tokenSecret string, locationIQAPIKey string) *Config {
	return &Config{dbClient, storeClient, tokenSecret, clock{}, locationIQAPIKey}
}
