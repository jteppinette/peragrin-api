package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

type Config struct {
	Client      *sqlx.DB
	TokenSecret string
}

func Init(client *sqlx.DB, tokenSecret string) *Config {
	return &Config{client, tokenSecret}
}
