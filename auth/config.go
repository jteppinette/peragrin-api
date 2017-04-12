package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

type Config struct {
	Client *sqlx.DB
}

func Init(client *sqlx.DB) *Config {
	return &Config{client}
}
