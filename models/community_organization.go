package models

import (
	"github.com/jmoiron/sqlx"
)

// CommunityOrganization represents the relationship between organizations and communities.
type CommunityOrganization struct {
	OrganizationID  int
	CommunityID     int
	IsAdministrator bool
}

// Create inserts a new community organization relationship in the database.
func (co *CommunityOrganization) Create(client *sqlx.DB) error {
	return client.Get(co, `
		INSERT INTO CommunityOrganization (organizationID, communityID, isAdministrator)
		VALUES ($1, $2, $3)
		RETURNING *;
	`, co.OrganizationID, co.CommunityID, co.IsAdministrator)
}

// Delete removes a community organization relationship from the database.
func (co *CommunityOrganization) Delete(client *sqlx.DB) error {
	if _, err := client.Exec("DELETE FROM CommunityOrganization WHERE organizationID = $1 AND communityID = $2;", co.OrganizationID, co.CommunityID); err != nil {
		return err
	}
	return nil
}
