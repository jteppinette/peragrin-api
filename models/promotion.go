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
	MembershipID   *int                `json:"membershipID,omitempty"`

	// Redemptions is the number of times this promotion has been redeemed.
	Redemptions int `json:"redemptions,omitempty"`
}

// Save creates or updates a promotion in the database based on the existence of an id.
func (p *Promotion) Save(client *sqlx.DB) error {
	if p.ID != 0 {
		return client.Get(p, "UPDATE Promotion SET name = $2, description = $3, exclusions = $4, expiration = $5, isSingleUse = $6, membershipID = $7 WHERE id = $1 RETURNING *;", p.ID, p.Name, p.Description, p.Exclusions, p.Expiration, p.IsSingleUse, p.MembershipID)
	}
	return client.Get(p, "INSERT INTO Promotion (organizationID, name, description, exclusions, expiration, isSingleUse, membershipID) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;", p.OrganizationID, p.Name, p.Description, p.Exclusions, p.Expiration, p.IsSingleUse, p.MembershipID)
}

// DeletePromotion removes a promotion from the database.
func DeletePromotion(id int, client *sqlx.DB) error {
	_, err := client.Exec("DELETE FROM Promotion WHERE id = $1;", id)
	if err != nil {
		return err
	}
	return nil
}

// GetPromotionsByOrganization returns all promotions for a given organization.
func GetPromotionsByOrganization(organizationID int, client *sqlx.DB) (Promotions, error) {
	promotions := Promotions{}
	if err := client.Select(&promotions, "SELECT Promotion.*, COUNT(AccountPromotion) AS redemptions FROM Promotion LEFT OUTER JOIN AccountPromotion ON (Promotion.id = AccountPromotion.promotionID) WHERE Promotion.organizationID = $1 GROUP BY Promotion.id;", organizationID); err != nil {
		return nil, err
	}
	return promotions, nil
}

// GetPromotionByID returns the promotion with the provided ID.
func GetPromotionByID(id int, client *sqlx.DB) (*Promotion, error) {
	promotion := &Promotion{}
	if err := client.Get(promotion, "SELECT * FROM Promotion WHERE id = $1;", id); err != nil {
		return nil, err
	}
	return promotion, nil
}

// GetPromotionsByID returns the promotions that have an id in the provided list of ids.
func GetPromotionsByID(ids []int, client *sqlx.DB) (Promotions, error) {
	promotions := Promotions{}
	query, args, err := sqlx.In("SELECT * FROM Promotion WHERE id IN (?);", ids)
	if err != nil {
		return nil, err
	}
	if err := client.Select(&promotions, client.Rebind(query), args...); err != nil {
		return nil, err
	}
	return promotions, nil
}
