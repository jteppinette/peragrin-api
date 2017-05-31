package models

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/peragrin/api/common"
)

// Promotions is a list of promotion objects.
type Promotions []Promotion

// Promotion represents a single promotion that is created by a community.
type Promotion struct {
	ID             int                 `json:"id"`
	OrganizationID int                 `json:"organizationID"`
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	Exclusions     string              `json:"exclusions"`
	Expiration     common.JSONNullTime `json:"expiration"`
	IsSingleUse    bool                `json:"isSingleUse"`
}

// Save creates or updates a promotion in the database based on the existence of an id.
func (p *Promotion) Save(client *sqlx.DB) error {
	if p.ID != 0 {
		return client.Get(p, "UPDATE Promotion SET name = $2, description = $3, exclusions = $4, expiration = $5, isSingleUse = $6 WHERE id = $1 RETURNING *;", p.ID, p.Name, p.Description, p.Exclusions, p.Expiration, p.IsSingleUse)
	}
	return client.Get(p, "INSERT INTO Promotion (organizationID, name, description, exclusions, expiration, isSingleUse) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;", p.OrganizationID, p.Name, p.Description, p.Exclusions, p.Expiration, p.IsSingleUse)
}

// GetPromotionsByOrganization returns all promotions for a given organization.
func GetPromotionsByOrganization(organizationID int, client *sqlx.DB) (Promotions, error) {
	promotions := Promotions{}
	if err := client.Select(&promotions, "SELECT * FROM Promotion WHERE OrganizationID = $1", organizationID); err != nil {
		return nil, err
	}
	return promotions, nil
}