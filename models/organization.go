package models

import (
	geo "github.com/codingsince1985/geo-golang"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Organizations is a list of organization structs.
type Organizations []Organization

// Organization represents an organization that is registered in
// the Peragrin system.
type Organization struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// SetGeo does a reverse geo-code lookup to turn an address into coordinates.
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

// Save creates or updates an organization based on the existence of an id.
func (o *Organization) Save(client *sqlx.DB) error {
	if o.ID != 0 {
		return client.Get(o, "UPDATE Organization SET name = $2, address = $3, longitude = $7, latitude = $8 WHERE id = $1 RETURNING *;", o.ID, o.Name, o.Address, o.Longitude, o.Latitude)
	}
	return client.Get(o, "INSERT INTO Organization (name, address, longitude, latitude) VALUES ($1, $2, $3, $4) RETURNING *;", o.Name, o.Address, o.Longitude, o.Latitude)
}

// AddOperator creates a relationship that defines a given account as
// an operator of this organization.
func (o *Organization) AddOperator(accountID int, client *sqlx.DB) error {
	operator := Operator{AccountID: accountID, OrganizationID: o.ID}
	return operator.Save(client)
}

// AddMembership creates an objects that describes the relationship between
// this organization and the provided community. If `isAdministrator` is true, then
// this organization will be an administrator over the provided community.
func (o *Organization) AddMembership(communityID int, isAdministrator bool, client *sqlx.DB) error {
	membership := Membership{OrganizationID: o.ID, CommunityID: communityID, IsAdministrator: isAdministrator}
	return membership.Save(client)
}

// ListOrganizations returns all organizations in the database.
func ListOrganizations(client *sqlx.DB) (Organizations, error) {
	organizations := Organizations{}
	if err := client.Select(&organizations, "SELECT * FROM Organization;"); err != nil {
		return nil, err
	}
	return organizations, nil
}

// ListOrganizationsByCommunityID returns all organizations in a given community.
// TODO: This SQL statement needs to be converted into a join across the M2M
// intermediary relation.
func ListOrganizationsByCommunityID(id int, client *sqlx.DB) (Organizations, error) {
	organizations := Organizations{}
	if err := client.Select(&organizations, "SELECT Organization.* FROM Organization INNER JOIN Membership ON (Organization.ID = Membership.OrganizationID) WHERE communityid = $1;", id); err != nil {
		return nil, err
	}
	return organizations, nil
}

// GetOrganizationByID returns the organization with the given id.
func GetOrganizationByID(id int, client *sqlx.DB) (*Organization, error) {
	o := &Organization{}
	if err := client.Get(o, "SELECT * FROM Organization WHERE id = $1;", id); err != nil {
		return nil, err
	}
	return o, nil
}
