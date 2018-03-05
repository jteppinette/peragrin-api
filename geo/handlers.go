package geo

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/jteppinette/peragrin-api/models"
	"github.com/jteppinette/peragrin-api/service"
)

// LookupHandler returns the longitude and latitude given an address.
func (c *Config) LookupHandler(r *http.Request) *service.Response {
	address := models.Address{}
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	lon, lat, err := address.Geocode(c.LocationIQAPIKey)
	if err != nil {
		log.WithFields(log.Fields{
			"street": address.Street, "city": address.City, "state": address.State, "country": address.Country, "zip": address.Zip,
			"error": err.Error(),
			"id":    r.Header.Get("X-Request-ID"),
		}).Info(errGeocode.Error())
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}{lon, lat})
}
