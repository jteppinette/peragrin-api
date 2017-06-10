package organizations

import (
	"github.com/jmoiron/sqlx"
	minio "github.com/minio/minio-go"
)

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	DBClient    *sqlx.DB
	StoreClient *minio.Client
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(dbClient *sqlx.DB, storeClient *minio.Client) *Config {
	return &Config{dbClient, storeClient}
}
