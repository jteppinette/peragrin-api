package models

import (
	"fmt"
	"net/url"
	"strconv"

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

// GetCommunityByID returns the requested community.
func GetCommunityByID(id int, client *sqlx.DB) (Community, error) {
	community := Community{}
	if err := client.Get(&community, "SELECT * FROM Community WHERE id = $1;", id); err != nil {
		return community, err
	}
	return community, nil
}

// GetCommunityByMembershipID returns the requested community.
func GetCommunityByMembershipID(id int, client *sqlx.DB) (Community, error) {
	community := Community{}
	if err := client.Get(&community, "SELECT Community.* FROM Community INNER JOIN Membership ON (Community.ID = Membership.CommunityID) WHERE Membership.ID = $1;", id); err != nil {
		return community, err
	}
	return community, nil
}

// GetCommunitiesByAccount returns all communities that are connected by the given account.
// This function will also return only those communities that are connected via an isAdministrator
// link, if the isAdministrator field is in the provided query.
func GetCommunitiesByAccount(accountID int, query url.Values, client *sqlx.DB) (Communities, error) {
	communities := Communities{}
	const format string = `
		SELECT DISTINCT ON (Community.id) Community.*, CommunityOrganization.isAdministrator FROM Community
		INNER JOIN CommunityOrganization ON (Community.id = CommunityOrganization.communityID)
		INNER JOIN AccountOrganization ON (CommunityOrganization.organizationID = AccountOrganization.organizationID)
		WHERE AccountOrganization.accountID = $1 %s ORDER BY Community.id, CommunityOrganization.isAdministrator DESC;
	`
	if b, err := strconv.ParseBool(query.Get("isAdministrator")); err != nil {
		if err := client.Select(&communities, fmt.Sprintf(format, ""), accountID); err != nil {
			return nil, err
		}
	} else {
		if err := client.Select(&communities, fmt.Sprintf(format, "AND CommunityOrganization.isAdministrator = $2"), accountID, b); err != nil {
			return nil, err
		}

	}
	return communities, nil
}
