package models

import (
	"github.com/jmoiron/sqlx"
)

// AccountOrganization represents an entity that can login
// into the Peragrin system and manage a organization or organizations.
type AccountOrganization struct {
	AccountID      int
	OrganizationID int
}

// Create persists the given operator in the database.
func (ao *AccountOrganization) Create(client *sqlx.DB) error {
	return client.Get(ao, "INSERT INTO AccountOrganization (accountID, organizationID) VALUES ($1, $2) RETURNING *;", ao.AccountID, ao.OrganizationID)
}
