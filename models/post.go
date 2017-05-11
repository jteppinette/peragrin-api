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
		return client.Get(p, "UPDATE Post SET content = $2, organizationID = $3 WHERE id = $1 RETURNING *;", p.ID, p.Content, p.OrganizationID)
	} else {
		return client.Get(p, "INSERT INTO Post (content, organizationID) VALUES ($1, $2) RETURNING *;", p.Content, p.OrganizationID)
	}
}

func ListPostsByCommunityID(id int, client *sqlx.DB) (Posts, error) {
	posts := Posts{}
	// TODO: Join the necessary tables to retreive all posts in all organizations
	// with the given community id.
	if err := client.Select(&posts, "SELECT * FROM Post ORDER BY createdAt DESC;"); err != nil {
		return nil, err
	}
	return posts, nil
}
