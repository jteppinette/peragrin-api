package models

import (
	"github.com/jmoiron/sqlx"
)

// Membership represents a level of membership that
// a patron could possess with a community.
type Membership struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Save creates or updates a membership in the database based on the existence of an id.
func (m *Membership) Save(communityID int, client *sqlx.DB) error {
	if m.ID != 0 {
		return client.Get(m, "UPDATE Membership SET communityID = $2, name = $3, description = $4 WHERE id = $1 RETURNING id, name, description;", m.ID, communityID, m.Name, m.Description)
	}
	return client.Get(m, "INSERT INTO Membership (communityID, name, description) VALUES ($1, $2, $3) RETURNING id, name, description;", communityID, m.Name, m.Description)
}

// GetMembershipsByCommunity returns all of a communities' memberships.
func GetMembershipsByCommunity(communityID int, client *sqlx.DB) ([]Membership, error) {
	memberships := []Membership{}
	if err := client.Select(&memberships, "SELECT id, name, description FROM Membership WHERE communityID = $1;", communityID); err != nil {
		return nil, err
	}
	return memberships, nil
}
