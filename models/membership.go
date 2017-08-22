package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Membership represents a level of membership that
// a patron could possess with a community.
type Membership struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CommunityID int     `json:"communityID,omitempty"`
	Price       float64 `json:"price"`

	// Expiration is used to define the time left for a provided account membership.
	// This information is only useful when in the context of an account.
	Expiration time.Time `json:"expiration,omitempty"`
}

// Create adds a new membership to the database.
func (m *Membership) Create(communityID int, client *sqlx.DB) error {
	return client.Get(m, "INSERT INTO Membership (communityID, name, description, price) VALUES ($1, $2, $3, $4) RETURNING id, name, description, price;", communityID, m.Name, m.Description, m.Price)
}

// Update updates a membership row in the database.
func (m *Membership) Update(client *sqlx.DB) error {
	return client.Get(m, "UPDATE Membership SET name = $2, description = $3, price = $4 WHERE id = $1 RETURNING id, name, description, price;", m.ID, m.Name, m.Description, m.Price)
}

// GetMembershipsByCommunity returns all of a communities' memberships.
func GetMembershipsByCommunity(communityID int, client *sqlx.DB) ([]Membership, error) {
	memberships := []Membership{}
	if err := client.Select(&memberships, "SELECT id, name, description, price FROM Membership WHERE communityID = $1 ORDER BY price;", communityID); err != nil {
		return nil, err
	}
	return memberships, nil
}

// GetMembershipByID returns the requested membership.
func GetMembershipByID(membershipID int, client *sqlx.DB) (*Membership, error) {
	m := &Membership{}
	if err := client.Get(m, "SELECT id, name, description, price FROM Membership WHERE id = $1;", membershipID); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return m, nil
}

// GetMembershipsByAccount returns all memberships that an account has.
func GetMembershipsByAccount(accountID int, client *sqlx.DB) ([]Membership, error) {
	memberships := []Membership{}
	if err := client.Select(&memberships, `
		SELECT Membership.id, Membership.name, Membership.description, Membership.price, Membership.communityID, AccountMembership.expiration FROM Membership
		INNER JOIN AccountMembership ON (Membership.id = AccountMembership.membershipID)
		WHERE AccountMembership.accountID = $1 ORDER BY Membership.price;
	`, accountID); err != nil {
		return nil, err
	}
	return memberships, nil
}

// DeleteMembership removes a membership from the database.
func DeleteMembership(id int, client *sqlx.DB) error {
	_, err := client.Exec("DELETE FROM Membership WHERE id = $1;", id)
	if err != nil {
		return err
	}
	return nil
}
