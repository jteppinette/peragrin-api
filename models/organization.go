package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

// Organizations is a list of organization structs.
type Organizations []Organization

// Organization represents an organization that is registered in
// the Peragrin system.
type Organization struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Address
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Address represents a physical location in the world.
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"Country"`
	Zip     string `json:"zip"`
}

func (a Address) geocode(key string) (float64, float64, error) {
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

// SetGeo does a reverse geo-code lookup to turn an address into coordinates.
func (o *Organization) SetGeo(key string) error {
	if o.Street == "" || o.City == "" || o.State == "" || o.Country == "" || o.Zip == "" {
		return errAddressRequired
	}

	var err error
	if o.Lon, o.Lat, err = o.geocode(key); err != nil {
		return err
	}

	return nil
}

// Save creates or updates an organization based on the existence of an id.
func (o *Organization) Save(client *sqlx.DB) error {
	if o.ID != 0 {
		return client.Get(o, "UPDATE Organization SET name = $2, street = $3, city = $4, state = $5, country = $6, zip = $7, lon = $8, lat = $9 WHERE id = $1 RETURNING *;", o.ID, o.Name, o.Street, o.City, o.State, o.Country, o.Zip, o.Lon, o.Lat)
	}
	return client.Get(o, "INSERT INTO Organization (name, street, city, state, country, zip, lon, lat) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;", o.Name, o.Street, o.City, o.State, o.Country, o.Zip, o.Lon, o.Lat)
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
