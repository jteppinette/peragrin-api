package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Posts []Post

type Post struct {
	ID             int       `json:"id"`
	Content        string    `json:"content"`
	OrganizationID int       `json:"organizationID"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (p *Post) Save(client *sqlx.DB) error {
	if p.ID != 0 {
		return client.Get(p, "UPDATE posts SET content = $2, organizationID = $3 WHERE id = $1 RETURNING *;", p.ID, p.Content, p.OrganizationID)
	} else {
		return client.Get(p, "INSERT INTO posts (content, organizationID) VALUES ($1, $2) RETURNING *;", p.Content, p.OrganizationID)
	}
}