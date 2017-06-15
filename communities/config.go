package communities

import (
	"github.com/jmoiron/sqlx"
	minio "github.com/minio/minio-go"
)

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	DBClient         *sqlx.DB
	StoreClient      *minio.Client
	LocationIQAPIKey string
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(dbClient *sqlx.DB, storeClient *minio.Client, locationIQAPIKey string) *Config {
	return &Config{dbClient, storeClient, locationIQAPIKey}
}
