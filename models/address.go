package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Address represents a physical location in the world.
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
}

func (a Address) Geocode(key string) (float64, float64, error) {
	r, err := http.Get(fmt.Sprintf("http://locationiq.org/v1/search.php?key=%s&format=json&limit=1&street=%s&city=%s&state=%s&country=%s&postalcode=%s", key, a.Street, a.City, a.State, a.Country, a.Zip))
	if err != nil {
		return 0, 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("expected response to be HTTP 200, received %s", r.Status)
	}

	codes := []struct {
		Lon string `json:"lon"`
		Lat string `json:"lat"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&codes); err != nil {
		return 0, 0, err
	}
	for _, code := range codes {
		lon, _ := strconv.ParseFloat(code.Lon, 64)
		lat, _ := strconv.ParseFloat(code.Lat, 64)
		return lon, lat, nil
	}
	return 0, 0, errGeocodeNotFound
}
