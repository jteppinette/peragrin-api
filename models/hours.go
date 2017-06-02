package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Hours represents a single days hours of operation.
type Hours []struct {
	Weekday time.Weekday `json:"weekday"`
	Start   int          `json:"start"`
	Close   int          `json:"close"`
}

// Set replaces an oragnizations hours of operation.
func (h Hours) Set(organizationID int, client *sqlx.DB) error {
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

	_, err = tx.Exec("DELETE FROM Hours WHERE organizationID = $1", organizationID)
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

	_, err = tx.Exec(client.Rebind(statement[0:len(statement)-1]), args...)
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
