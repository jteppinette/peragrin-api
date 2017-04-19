package auth

import (
	geo "github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/mapbox"
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

type Config struct {
	Client       *sqlx.DB
	TokenSecret  string
	Clock        timer
	Geo          geo.Geocoder
	MapboxAPIKey string
}

func Init(client *sqlx.DB, tokenSecret string, mapboxAPIKey string) *Config {
	return &Config{client, tokenSecret, clock{}, mapbox.Geocoder(mapboxAPIKey), mapboxAPIKey}
}
