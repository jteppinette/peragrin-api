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
			id		SERIAL PRIMARY KEY,
			name	varchar(80) NOT NULL UNIQUE,
			street	varchar(160) NOT NULL,
			city	varchar(30) NOT NULL,
			state	varchar(30) NOT NULL,
			country	varchar(40) NOT NULL,
			zip		varchar(20) NOT NULL,
			lon		float,
			lat		float,
			email varchar(60) NOT NULL,
			phone varchar(20) NOT NULL,
			website varchar(60) NOT NULL,
			category varchar(30)
		);

		CREATE TABLE IF NOT EXISTS Hours (
			organizationID integer REFERENCES Organization ON DELETE CASCADE,
			weekday int,
			start int,
			close int,
			PRIMARY KEY (organizationID, weekday, start, close)
		);

		CREATE TABLE IF NOT EXISTS Promotion (
			id				SERIAL PRIMARY KEY,
			organizationID	integer REFERENCES Organization ON DELETE CASCADE,
			name			varchar(80) NOT NULL,
			description		text,
			exclusions		text,
			expiration		date,
			isSingleUse		bool
		);
		CREATE UNIQUE INDEX ON Promotion (organizationID, name);

		CREATE TABLE IF NOT EXISTS Community (
			id			SERIAL PRIMARY KEY,
			name		varchar(80) NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS GeoJSONOverlay (
			name			varchar(40),
			communityID		integer REFERENCES Community ON DELETE CASCADE,
			data			jsonb,
			style			jsonb,
			PRIMARY KEY (communityID, name)
		);

		CREATE TABLE IF NOT EXISTS Post (
			id				SERIAL PRIMARY KEY,
			content			text,
			organizationID	integer NOT NULL REFERENCES Organization ON DELETE CASCADE,
			createdAt		timestamp DEFAULT current_timestamp
		);

		CREATE TABLE IF NOT EXISTS AccountOrganization (
			accountID			int REFERENCES Account ON DELETE CASCADE,
			organizationID		int REFERENCES Organization ON DELETE CASCADE,
			PRIMARY KEY (accountID, organizationID)
		);

		CREATE TABLE IF NOT EXISTS CommunityOrganization (
			organizationID	int REFERENCES Organization ON DELETE CASCADE,
			communityID		int REFERENCES Community ON DELETE CASCADE,
			isAdministrator bool,
			PRIMARY KEY (organizationID, communityID)
		);
	`
	if _, err := client.Exec(schema); err != nil {
		return err
	}
	return nil
}
