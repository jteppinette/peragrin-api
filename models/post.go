package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Posts is a list of post objects.
type Posts []Post

// Post represents a single social media post by an organization.
type Post struct {
	ID             int       `json:"id"`
	Content        string    `json:"content"`
	OrganizationID int       `json:"organizationID"`
	CreatedAt      time.Time `json:"createdAt"`
}

// Save creates or updates a post in the database based on the existence of an id.
func (p *Post) Save(client *sqlx.DB) error {
	if p.ID != 0 {
		return client.Get(p, "UPDATE Post SET content = $2, organizationID = $3 WHERE id = $1 RETURNING *;", p.ID, p.Content, p.OrganizationID)
	}
	return client.Get(p, "INSERT INTO Post (content, organizationID) VALUES ($1, $2) RETURNING *;", p.Content, p.OrganizationID)
}

// GetPostsByCommunity returns all posts created by organizations that are a
// member of the provided community.
func GetPostsByCommunity(communityID int, client *sqlx.DB) (Posts, error) {
	posts := Posts{}
	if err := client.Select(&posts, "SELECT Post.* FROM Post INNER JOIN CommunityOrganization ON (Post.organizationID = CommunityOrganization.organizationID) WHERE communityID = $1 ORDER BY createdAt DESC;", communityID); err != nil {
		return nil, err
	}
	return posts, nil
}
