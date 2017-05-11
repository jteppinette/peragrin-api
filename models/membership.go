package models

import (
	"github.com/jmoiron/sqlx"
)

// Memberships is a slice of membership structs.
type Memberships []Membership

// Membership represents the relationship between organizations and communities.
type Membership struct {
	OrganizationID  int  `json:"organizationID"`
	CommunityID     int  `json:"communityID"`
	IsAdministrator bool `json:"isAdministrator"`
}

// Save creates the given membership in the database.
func (m *Membership) Save(client *sqlx.DB) error {
	return client.Get(m, "INSERT INTO Membership (organizationID, communityID, isAdministrator) VALUES ($1, $2, $3) RETURNING *;", m.OrganizationID, m.CommunityID, m.IsAdministrator)
}

// ListMembershipsByOrganizationID returns all memberships of the provided organization.
func ListMembershipsByOrganizationID(organizationID int, client *sqlx.DB) (Memberships, error) {
	memberships := Memberships{}
	if err := client.Select(&memberships, "SELECT * FROM Membership WHERE organizationID = $1;", organizationID); err != nil {
		return nil, err
	}
	return memberships, nil
}

// ListMembershipsByCommunityID returns all of a communities' memberships.
func ListMembershipsByCommunityID(communityID int, client *sqlx.DB) (Memberships, error) {
	memberships := Memberships{}
	if err := client.Select(&memberships, "SELECT * FROM Membership WHERE communityID = $1;", communityID); err != nil {
		return nil, err
	}
	return memberships, nil
}
