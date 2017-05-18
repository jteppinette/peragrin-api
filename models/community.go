package models

import (
	"github.com/jmoiron/sqlx"
)

// Communities represents a list of community objects.
type Communities []Community

// Community is a geographic area that supports joint interaction
// between organizations and patrons.
type Community struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// IsAdministrator is only populated when this community
	// is in the context of an organization.
	IsAdministrator *bool `json:"isAdministrator,omitempty"`
}

// Create persists the provided community in the database, and it creates
// the relationship to the provided organization. This will be an administrative
// relationship.
func (c *Community) Create(organizationID int, client *sqlx.DB) error {
	if err := client.Get(c, "INSERT INTO Community (name) VALUES ($1) RETURNING *;", c.Name); err != nil {
		return err
	}
	co := CommunityOrganization{OrganizationID: organizationID, CommunityID: c.ID, IsAdministrator: true}
	if err := co.Create(client); err != nil {
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
