package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	client *sqlx.DB
)

func Client(host, user, password, dbname string) (*sqlx.DB, error) {
	if client != nil {
		return client, nil
	}
	var err error
	client, err = sqlx.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname))
	if err != nil {
		return nil, err
	}
	return client, nil
}
