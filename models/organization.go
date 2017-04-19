package models

import (
	geo "github.com/codingsince1985/geo-golang"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Organizations []Organization

type Organization struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Leader      bool    `json:"leader"`
	Enabled     bool    `json:"enabled"`
	CommunityID int     `json:"communityID"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}

func (o *Organization) SetGeo(geocoder geo.Geocoder) error {
	if o.Address == "" {
		return errAddressRequired
	}
	location, err := geocoder.Geocode(o.Address)
	if err != nil {
		return errors.Wrap(err, errGeo.Error())
	}
	if location == nil || location.Lat == 0 || location.Lng == 0 {
		return errGeo
	}
	o.Longitude = location.Lng
	o.Latitude = location.Lat
	return nil
}

func (o *Organization) Save(client *sqlx.DB) error {
	if o.ID != 0 {
		return client.Get(o, "UPDATE organizations SET name = $2, address = $3, leader = $4, enabled = $5, communityID = $6, longitude = $7, latitude = %8 WHERE id = $1 RETURNING *;", o.ID, o.Name, o.Address, o.Leader, o.Enabled, o.CommunityID, o.Longitude, o.Latitude)
	} else {
		return client.Get(o, "INSERT INTO organizations (name, address, leader, enabled, communityID, longitude, latitude) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;", o.Name, o.Address, o.Leader, o.Enabled, o.CommunityID, o.Longitude, o.Latitude)
	}
}

func ListOrganizations(client *sqlx.DB) (Organizations, error) {
	organizations := Organizations{}
	if err := client.Select(&organizations, "SELECT * FROM organizations;"); err != nil {
		return nil, err
	}
	return organizations, nil
}

func ListOrganizationsByCommunityID(id int, client *sqlx.DB) (Organizations, error) {
	organizations := Organizations{}
	if err := client.Select(&organizations, "SELECT * FROM organizations WHERE communityID = $1;", id); err != nil {
		return nil, err
	}
	return organizations, nil
}

func GetOrganizationByID(id int, client *sqlx.DB) (*Organization, error) {
	o := &Organization{}
	if err := client.Get(o, "SELECT * FROM organizations WHERE id = $1;", id); err != nil {
		return nil, err
	}
	return o, nil
}
