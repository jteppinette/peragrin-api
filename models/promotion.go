package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	Communities    pq.Int64Array       `json:"communities"`

	// Redemptions is the number of times this promotion has been redeemed.
	Redemptions int `json:"redemptions,omitempty"`
}

// Save creates or updates a promotion in the database based on the existence of an id.
func (p *Promotion) Save(client *sqlx.DB) error {
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

	if p.ID != 0 {
		err = tx.Get(p, "UPDATE Promotion SET name = $2, description = $3, exclusions = $4, expiration = $5, isSingleUse = $6 WHERE id = $1 RETURNING *;", p.ID, p.Name, p.Description, p.Exclusions, p.Expiration, p.IsSingleUse)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM CommunityPromotion WHERE promotionID = $1;", p.ID)
		if err != nil {
			return err
		}
	} else {
		err = tx.Get(p, "INSERT INTO Promotion (organizationID, name, description, exclusions, expiration, isSingleUse) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;", p.OrganizationID, p.Name, p.Description, p.Exclusions, p.Expiration, p.IsSingleUse)
		if err != nil {
			return err
		}
	}

	if len(p.Communities) == 0 {
		return nil
	}

	statement := "INSERT INTO CommunityPromotion (communityID, promotionID) VALUES "
	args := make([]interface{}, len(p.Communities)*2)

	for i, communityID := range p.Communities {
		statement = statement + "(?, ?),"
		set := i * 2
		args[set+0] = communityID
		args[set+1] = p.ID
	}

	_, err = tx.Exec(sqlx.Rebind(sqlx.BindType("postgres"), statement[0:len(statement)-1]), args...)
	if err != nil {
		return err
	}

	return nil
}

// DeletePromotion removes a promotion from the database.
func DeletePromotion(id int, client *sqlx.DB) error {
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

	_, err = tx.Exec("DELETE FROM Promotion WHERE id = $1;", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM CommunityPromotion WHERE promotionID = $1;", id)
	if err != nil {
		return err
	}

	return nil
}

// GetPromotionsByOrganization returns all promotions for a given organization.
func GetPromotionsByOrganization(organizationID int, client *sqlx.DB) (Promotions, error) {
	promotions := Promotions{}
	if err := client.Select(&promotions, `
		SELECT Promotion.*, COUNT(AccountPromotion) AS redemptions, (
			SELECT ARRAY(SELECT communityID FROM CommunityPromotion WHERE promotionID = Promotion.id)
		) as communities
		FROM Promotion LEFT OUTER JOIN AccountPromotion ON (Promotion.id = AccountPromotion.promotionID)
		WHERE Promotion.organizationID = $1 GROUP BY Promotion.id;
	`, organizationID); err != nil {
		return nil, err
	}
	return promotions, nil
}

// GetPromotionsByID returns the promotions that have an id in the provided list of ids.
func GetPromotionsByID(ids []int, client *sqlx.DB) (Promotions, error) {
	promotions := Promotions{}
	query, args, err := sqlx.In(`
		SELECT Promotion.*, (
			SELECT ARRAY(SELECT communityID FROM CommunityPromotion WHERE promotionID = Promotion.id)
		) as communities
		FROM Promotion WHERE Promotion.id IN (?);
	`, ids)
	if err != nil {
		return nil, err
	}
	if err := client.Select(&promotions, client.Rebind(query), args...); err != nil {
		return nil, err
	}

	if err := client.Select(&promotions, client.Rebind(query), args...); err != nil {
		return nil, err
	}
	return promotions, nil
}
