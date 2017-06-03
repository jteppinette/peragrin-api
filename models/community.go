package models

import (
	"github.com/jmoiron/sqlx"
)

// Communities represents a list of community objects.
type Communities []Community

// Community is a geographic area that supports joint interaction
// between organizations and patrons.
type Community struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Lon  float64 `json:"lon"`
	Lat  float64 `json:"lat"`
	Zoom int     `json:"zoom"`

	// IsAdministrator is only populated when this community
	// is in the context of an organization.
	IsAdministrator *bool `json:"isAdministrator,omitempty"`
}

// Create persists the provided community in the database, and it creates
// the relationship to the provided organization. This will be an administrative
// relationship.
func (c *Community) Create(organizationID int, client *sqlx.DB) error {
	tx, err := client.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = tx.Get(c, "INSERT INTO Community (name, lon, lat, zoom) VALUES ($1, $2, $3, $4) RETURNING *;", c.Name, c.Lon, c.Lat, c.Zoom)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO CommunityOrganization (organizationID, communityID, isAdministrator) VALUES ($1, $2, $3)", organizationID, c.ID, true)
	if err != nil {
		return err
	}

	return nil
}

// GetCommunities returns all communities in the database.
func GetCommunities(client *sqlx.DB) (Communities, error) {
	communities := Communities{}
	if err := client.Select(&communities, "SELECT * FROM Community;"); err != nil {
		return nil, err
	}

	return communities, nil
}

// GetCommunitiesByOrganization returns all communities with a relationship
// to the provided organization.
func GetCommunitiesByOrganization(organizationID int, client *sqlx.DB) (Communities, error) {
	communities := Communities{}
	if err := client.Select(&communities, "SELECT Community.*, CommunityOrganization.isAdministrator FROM Community INNER JOIN CommunityOrganization ON (Community.id = CommunityOrganization.communityID) WHERE organizationID = $1", organizationID); err != nil {
		return nil, err
	}
	return communities, nil
}
