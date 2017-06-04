package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// AccountPromotion represents the relationship between organizations and communities.
type AccountPromotion struct {
	AccountID   int       `json:"-"`
	PromotionID int       `json:"-"`
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
