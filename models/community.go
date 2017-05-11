package models

import (
	"github.com/jmoiron/sqlx"
)

type Communities []Community

type Community struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Community) Save(client *sqlx.DB) error {
	if c.ID != 0 {
		return client.Get(c, "UPDATE Community SET name = $2 WHERE id = $1 RETURNING *;", c.ID, c.Name)
	} else {
		return client.Get(c, "INSERT INTO Community (name) VALUES ($1) RETURNING *;", c.Name)
	}
}

func ListCommunities(client *sqlx.DB) (Communities, error) {
	communities := Communities{}
	if err := client.Select(&communities, "SELECT * FROM Community;"); err != nil {
		return nil, err
	}
	return communities, nil
}
