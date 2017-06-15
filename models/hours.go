package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Hour is a single day's schedule.
type Hour struct {
	Weekday time.Weekday `json:"weekday"`
	Start   int          `json:"start"`
	Close   int          `json:"close"`
}

// Hours represents the full week schedule.
type Hours []Hour

// Set replaces an organizations hours of operation.
func (h Hours) txSet(organizationID int, tx *sqlx.Tx) error {
	_, err := tx.Exec("DELETE FROM Hours WHERE organizationID = $1", organizationID)
	if err != nil {
		return err
	}

	statement := "INSERT INTO Hours (organizationID, weekday, start, close) VALUES "
	args := make([]interface{}, len(h)*4)

	for i, v := range h {
		statement = statement + "(?, ?, ?, ?),"
		set := i * 4
		args[set+0] = organizationID
		args[set+1] = v.Weekday
		args[set+2] = v.Start
		args[set+3] = v.Close
	}

	_, err = tx.Exec(sqlx.Rebind(sqlx.BindType("postgres"), statement[0:len(statement)-1]), args...)
	if err != nil {
		return err
	}

	return nil
}

// GetHoursByOrganization returns the set of hours for the given organization.
func GetHoursByOrganization(organizationID int, client *sqlx.DB) (Hours, error) {
	hours := Hours{}
	if err := client.Select(&hours, "SELECT Hours.weekday, Hours.start, Hours.close FROM Hours WHERE organizationID = $1;", organizationID); err != nil {
		return nil, err
	}
	return hours, nil
}
