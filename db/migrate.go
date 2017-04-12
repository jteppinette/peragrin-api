package db

import (
	"github.com/jmoiron/sqlx"
)

func Migrate(client *sqlx.DB) error {
	const schema = `
		CREATE TABLE IF NOT EXISTS users (
			id			SERIAL,
			username	varchar(40) NOT NULL UNIQUE,
			password	varchar(60) NOT NULL
		);
	`
	if _, err := client.Exec(schema); err != nil {
		return err
	}
	return nil
}
