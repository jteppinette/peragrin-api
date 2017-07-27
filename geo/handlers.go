package geo

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// LookupHandler returns the longitude and latitude given an address.
func (c *Config) LookupHandler(r *http.Request) *service.Response {
	address := models.Address{}
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	log.Info(address)

	lon, lat, err := address.Geocode(c.LocationIQAPIKey)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}{lon, lat})
}
