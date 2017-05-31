package models

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	minio "github.com/minio/minio-go"
)

// Communities represents a list of community objects.
type Communities []Community

// Community is a geographic area that supports joint interaction
// between organizations and patrons.
type Community struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// GeoJSON is used to send the presigned GeoJSON file link to the client.
	GeoJSON string `json:"geoJSON"`

	// IsAdministrator is only populated when this community
	// is in the context of an organization.
	IsAdministrator *bool `json:"isAdministrator,omitempty"`
}

// UploadGeoJSON puts a new object in the static store.
func (c *Community) UploadGeoJSON(reader io.Reader, client *minio.Client) error {
	_, err := client.PutObject(bucket, fmt.Sprintf("geojson/%s", strconv.Itoa(c.ID)), reader, "application/octet-stream")
	return err
}

// SetPresignedGeoJSONLink sets the Logo field with a presigned get request url.
func (c *Community) SetPresignedGeoJSONLink(client *minio.Client) error {
	object := fmt.Sprintf("geojson/%s", strconv.Itoa(c.ID))

	// If this object does not exist, then do not set the presigned link.
	// TODO: Store a reference to the object in the database to determine existence.
	_, err := client.StatObject(bucket, object)
	if err != nil {
		return nil
	}

	url, err := client.PresignedGetObject(bucket, object, time.Second*24*60*60, nil)
	if err != nil {
		return err
	}
	c.GeoJSON = url.String()
	return nil
}

// SetPresignedGeoJSONLinks sets the Logo field with a presgned get request url for each organization provided.
func (communities Communities) SetPresignedGeoJSONLinks(client *minio.Client) error {
	for i, c := range communities {
		if err := c.SetPresignedGeoJSONLink(client); err != nil {
			return err
		}
		communities[i] = c
	}
	return nil
}

// Create persists the provided community in the database, and it creates
// the relationship to the provided organization. This will be an administrative
// relationship.
func (c *Community) Create(organizationID int, client *sqlx.DB) error {
	if err := client.Get(c, "INSERT INTO Community (name) VALUES ($1) RETURNING *;", c.Name); err != nil {
		return err
	}
	co := CommunityOrganization{OrganizationID: organizationID, CommunityID: c.ID, IsAdministrator: true}
	if err := co.Create(client); err != nil {
		return err
	}
	return nil
}

// GetCommunities returns all communities in the database.
func GetCommunities(client *sqlx.DB) (Communities, error) {
	communities := Communities{}
	if err := client.Select(&communities, "SELECT * FROM Community;"); err != nil {
		return nil, err
	}
	return communities, nil
}

// GetCommunitiesByOrganization returns all communities with a relationship
// to the provided organization.
func GetCommunitiesByOrganization(organizationID int, client *sqlx.DB) (Communities, error) {
	communities := Communities{}
	if err := client.Select(&communities, "SELECT Community.*, CommunityOrganization.isAdministrator FROM Community INNER JOIN CommunityOrganization ON (Community.id = CommunityOrganization.communityID) WHERE organizationID = $1", organizationID); err != nil {
		return nil, err
	}
	return communities, nil
}
