package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

type Config struct {
	Client           *sqlx.DB
	TokenSecret      string
	Clock            timer
	LocationIQAPIKey string
}

func Init(client *sqlx.DB, tokenSecret string, locationIQAPIKey string) *Config {
	return &Config{client, tokenSecret, clock{}, locationIQAPIKey}
}
