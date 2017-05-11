package db

import (
	"github.com/jmoiron/sqlx"
)

// Migrate prepares the database.
func Migrate(client *sqlx.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS Account (
			id				SERIAL PRIMARY KEY,
			email			varchar(60) NOT NULL UNIQUE,
			password		varchar(60) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS Organization (
			id				SERIAL PRIMARY KEY,
			name			varchar(80) NOT NULL UNIQUE,
			address			varchar(80) NOT NULL,
			longitude		float,
			latitude		float
		);
		CREATE TABLE IF NOT EXISTS Operator (
			accountID			int REFERENCES Account ON DELETE CASCADE,
			organizationID		int REFERENCES Organization ON DELETE CASCADE,
			PRIMARY KEY (accountID, organizationID)
		);
		CREATE TABLE IF NOT EXISTS Community (
			id			SERIAL PRIMARY KEY,
			name		varchar(80) NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS Membership (
			organizationID	int REFERENCES Organization ON DELETE CASCADE,
			communityID		int REFERENCES Community ON DELETE CASCADE,
			isAdministrator bool,
			PRIMARY KEY (organizationID, communityID)
		);
		CREATE TABLE IF NOT EXISTS Post (
			id				SERIAL PRIMARY KEY,
			content			text,
			organizationID	integer NOT NULL REFERENCES Organization ON DELETE CASCADE,
			createdAt		timestamp DEFAULT current_timestamp
		);
	`
	if _, err := client.Exec(schema); err != nil {
		return err
	}
	return nil
}
