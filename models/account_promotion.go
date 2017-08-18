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
		Exists bool
	}{}

	promotion, err := GetPromotionByID(ap.PromotionID, client)
	if err != nil {
		return false, err
	}

	if promotion.MembershipID == nil {
		return true, nil
	}

	if err := client.Get(&result, `
		SELECT EXISTS(
			SELECT FROM Account
			INNER JOIN AccountMembership ON (Account.id = AccountMembership.accountiD)
			INNER JOIN Promotion ON (AccountMembership.membershipID = Promotion.membershipID)
			WHERE Promotion.id = $1 AND Account.id  = $2
		);
	`, ap.PromotionID, ap.AccountID); err != nil {
		return false, err
	}
	return result.Exists, nil
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
