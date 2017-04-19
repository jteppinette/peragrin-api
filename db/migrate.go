package db

import (
	"github.com/jmoiron/sqlx"
)

// Migrate prepares the database.
func Migrate(client *sqlx.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id				SERIAL,
			email			varchar(60) NOT NULL UNIQUE,
			password		varchar(60) NOT NULL,
			organizationID	integer
		);
		CREATE TABLE IF NOT EXISTS communities (
			id			SERIAL,
			name		varchar(80) NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS organizations (
			id				SERIAL,
			name			varchar(80) NOT NULL UNIQUE,
			address			varchar(80) NOT NULL,
			isLeader		boolean,
			communityID		integer,
			longitude		float,
			latitude		float
		);
	`
	if _, err := client.Exec(schema); err != nil {
		return err
	}
	return nil
}
