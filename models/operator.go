package models

import (
	"github.com/jmoiron/sqlx"
)

// Operators is a slice of operator structs.
type Operators []Operator

// Operator represents an entity that can login into the Peragrin system
// and manage a organization or organizations.
type Operator struct {
	AccountID      int `json:"accountID"`
	OrganizationID int `json:"organizationID"`
}

// Save creates the given operator in the database.
func (o *Operator) Save(client *sqlx.DB) error {
	return client.Get(o, "INSERT INTO Operator (accountID, organizationID) VALUES ($1, $2) RETURNING *;", o.AccountID, o.OrganizationID)
}

// ListOperatorsByAccountID returns all of an account's operators.
func ListOperatorsByAccountID(accountID int, client *sqlx.DB) (Operators, error) {
	operators := Operators{}
	if err := client.Select(&operators, "SELECT * FROM Operator WHERE accountID = $1;", accountID); err != nil {
		return nil, err
	}
	return operators, nil
}

// ListOperatorsByOrganizationID returns all operators of the provided organization.
func ListOperatorsByOrganizationID(organizationID int, client *sqlx.DB) (Operators, error) {
	operators := Operators{}
	if err := client.Select(&operators, "SELECT * FROM Operator WHERE organizationID = $1;", organizationID); err != nil {
		return nil, err
	}
	return operators, nil
}
