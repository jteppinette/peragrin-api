package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// AccountPromotion represents the relationship between organizations and communities.
type AccountPromotion struct {
	AccountID   int       `json:"accountID, omitempty"`
	PromotionID int       `json:"promotionID, omitempty"`
	ConsumedAt  time.Time `json:"consumedAt"`
}

// Create inserts a new account promotion relationship in the database.
func (ap *AccountPromotion) Create(client *sqlx.DB) error {
	return client.Get(ap, `
		INSERT INTO AccountPromotion (accountID, promotionID)
		VALUES ($1, $2)
		RETURNING *;
	`, ap.AccountID, ap.PromotionID)
}

// HasPermission determines if the provided account has access to redeem
// the provided promotion.
func (ap *AccountPromotion) HasPermission(client *sqlx.DB) (bool, error) {
	result := struct {
		Exists   bool
		Required bool
	}{}

	if err := client.Get(&result, `
		SELECT
			EXISTS(
				SELECT FROM AccountMembership
				LEFT OUTER JOIN Membership ON (AccountMembership.membershipID = Membership.id)
				LEFT OUTER JOIN CommunityPromotion ON (Membership.communityID = CommunityPromotion.communityID)
				WHERE CommunityPromotion.promotionID = $1 AND AccountMembership.accountID = $2
			) AS exists,
			EXISTS(SELECT FROM CommunityPromotion WHERE promotionID = $1) AS required;
	`, ap.PromotionID, ap.AccountID); err != nil {
		return false, err
	}
	return result.Exists || !result.Required, nil
}

// GetAccountsPromotionsByID returns all account promotion redemption events for the given account and promotion.
func GetAccountsPromotionsByID(accountID, promotionID int, client *sqlx.DB) ([]AccountPromotion, error) {
	result := []AccountPromotion{}
	if err := client.Select(&result, "SELECT * FROM AccountPromotion WHERE AccountID = $1 AND PromotionID = $2 ORDER BY consumedAt DESC;", accountID, promotionID); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAccountsPromotionsByAccount returns all account promotion redemption events for the given account.
func GetAccountsPromotionsByAccount(accountID int, client *sqlx.DB) ([]AccountPromotion, error) {
	result := []AccountPromotion{}
	if err := client.Select(&result, "SELECT * FROM AccountPromotion WHERE AccountID = $1 ORDER BY consumedAt DESC;", accountID); err != nil {
		return nil, err
	}
	return result, nil
}
