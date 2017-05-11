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
}

// Save creates or updates a community based on the existence of an id.
func (c *Community) Save(client *sqlx.DB) error {
	if c.ID != 0 {
		return client.Get(c, "UPDATE Community SET name = $2 WHERE id = $1 RETURNING *;", c.ID, c.Name)
	}
	return client.Get(c, "INSERT INTO Community (name) VALUES ($1) RETURNING *;", c.Name)
}

// AddMembership creates an objects that describes the relationship between
// the provided organization and this community. If `isAdministrator` is true, then
// the provided organization will be an administrator over this community.
func (c *Community) AddMembership(organizationID int, isAdministrator bool, client *sqlx.DB) error {
	membership := Membership{OrganizationID: organizationID, CommunityID: c.ID, IsAdministrator: isAdministrator}
	return membership.Save(client)
}

// ListCommunities returns all communities in the database.
func ListCommunities(client *sqlx.DB) (Communities, error) {
	communities := Communities{}
	if err := client.Select(&communities, "SELECT * FROM Community;"); err != nil {
		return nil, err
	}
	return communities, nil
}
