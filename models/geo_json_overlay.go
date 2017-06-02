package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

// GeoJSONOverlay is an object that represents a list of features
// each determined by parameters and a set of geographical coordinates.
// These sets of geographical coordinates are called "Features". These
// features can be styled independtly through the Style object.
type GeoJSONOverlay struct {
	Name  string         `json:"name"`
	Data  types.JSONText `json:"data"`
	Style types.JSONText `json:"style"`
}

// Create creates a new geo json overlay in the database.
func (o *GeoJSONOverlay) Create(communityID int, client *sqlx.DB) error {
	return client.Get(o, "INSERT INTO GeoJSONOverlay (name, communityID, data, style) VALUES ($1, $2, $3, $4) RETURNING name, data, style;", o.Name, communityID, o.Data, o.Style)
}

// GetGeoJSONOverlaysByCommunity returns the set of geo JSON overlays that
// are used by the given community.
func GetGeoJSONOverlaysByCommunity(communityID int, client *sqlx.DB) ([]GeoJSONOverlay, error) {
	geoJSONOverlays := []GeoJSONOverlay{}
	if err := client.Select(&geoJSONOverlays, "SELECT name, data, style FROM GeoJSONOverlay WHERE communityID = $1;", communityID); err != nil {
		return nil, err
	}
	return geoJSONOverlays, nil
}
