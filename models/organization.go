package models

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	minio "github.com/minio/minio-go"
)

const bucket = "peragrin"

// Organizations is a list of organization structs.
type Organizations []Organization

// Organization represents an organization that is registered in
// the Peragrin system. This can be both a community leader's
// and business leader's organization.
type Organization struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Address
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
	Category string  `json:"category"`

	// Logo is used to send the presigned Logo file link to the client.
	Logo string `json:"logo"`

	// Icon is used to send the presigned Icon file link to the client.
	Icon string `json:"icon"`

	// IsAdministrator is only populated when this organization
	// is in the context of a community.
	IsAdministrator *bool `json:"isAdministrator,omitempty"`
}

// SetGeo does a reverse geo-code lookup to turn an address into coordinates.
func (o *Organization) SetGeo(key string) error {
	if o.Street == "" || o.City == "" || o.State == "" || o.Country == "" || o.Zip == "" {
		return errAddressRequired
	}
	var err error
	if o.Lon, o.Lat, err = o.geocode(key); err != nil {
		return err
	}
	return nil
}

// UploadLogo puts a new object in the static store.
func (o *Organization) UploadLogo(reader io.Reader, client *minio.Client) error {
	_, err := client.PutObject(bucket, fmt.Sprintf("logos/%s", strconv.Itoa(o.ID)), reader, "application/octet-stream")
	return err
}

// SetPresignedLogoLink sets the Logo field with a presigned get request url.
func (o *Organization) SetPresignedLogoLink(client *minio.Client) error {
	object := fmt.Sprintf("logos/%s", strconv.Itoa(o.ID))

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
	o.Logo = url.String()
	return nil
}

// SetPresignedLogoLinks sets the Logo field with a presgned get request url for each organization provided.
func (organizations Organizations) SetPresignedLogoLinks(client *minio.Client) error {
	for i, o := range organizations {
		if err := o.SetPresignedLogoLink(client); err != nil {
			return err
		}
		organizations[i] = o
	}
	return nil
}

// UploadIcon puts a new object in the static store.
func (o *Organization) UploadIcon(reader io.Reader, client *minio.Client) error {
	_, err := client.PutObject(bucket, fmt.Sprintf("icons/%s", strconv.Itoa(o.ID)), reader, "application/octet-stream")
	return err
}

// SetPresignedIconLink sets the Icon field with a presigned get request url.
func (o *Organization) SetPresignedIconLink(client *minio.Client) error {
	object := fmt.Sprintf("icons/%s", strconv.Itoa(o.ID))

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
	o.Icon = url.String()
	return nil
}

// SetPresignedIconLinks sets the Icon field with a presgned get request url for each organization provided.
// Only administrators have custom icons.
func (organizations Organizations) SetPresignedIconLinks(client *minio.Client) error {
	for i, o := range organizations {
		if o.IsAdministrator == nil {
			continue
		}
		if err := o.SetPresignedIconLink(client); err != nil {
			return err
		}
		organizations[i] = o
	}
	return nil
}

// Create persists a new organization in the database and creates the
// account - organization relationship.
func (o *Organization) Create(accountID int, client *sqlx.DB) error {
	if err := client.Get(o, `
		INSERT INTO Organization (name, street, city, state, country, zip, lon, lat, email, phone, website, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING *;
	`, o.Name, o.Street, o.City, o.State, o.Country, o.Zip, o.Lon, o.Lat, o.Email, o.Phone, o.Website, o.Category); err != nil {
		return err
	}
	ao := AccountOrganization{AccountID: accountID, OrganizationID: o.ID}
	if err := ao.Create(client); err != nil {
		return err
	}
	return nil
}

// Update updates the fields of a given organization.
func (o *Organization) Update(client *sqlx.DB) error {
	if err := client.Get(o, `
		UPDATE Organization
		SET name = $2, street = $3, city = $4, state = $5, country = $6, zip = $7, lon = $8, lat = $9, email = $10, phone = $11, website = $12, category = $13
		WHERE id = $1
		RETURNING *;
	`, o.ID, o.Name, o.Street, o.City, o.State, o.Country, o.Zip, o.Lon, o.Lat, o.Email, o.Phone, o.Website, o.Category); err != nil {
		return err
	}
	return nil
}

// GetOrganizationsByCommunity returns all organizations that are a member of the given community.
func GetOrganizationsByCommunity(communityID int, client *sqlx.DB) (Organizations, error) {
	organizations := Organizations{}
	if err := client.Select(&organizations, "SELECT Organization.*, CommunityOrganization.isAdministrator FROM Organization INNER JOIN CommunityOrganization ON (Organization.id = CommunityOrganization.organizationID) WHERE communityID = $1;", communityID); err != nil {
		return nil, err
	}
	return organizations, nil
}

// GetOrganizationsByAccount returns all organizations that are operated by the given account.
func GetOrganizationsByAccount(accountID int, client *sqlx.DB) (Organizations, error) {
	organizations := Organizations{}
	if err := client.Select(&organizations, "SELECT Organization.* FROM Organization INNER JOIN AccountOrganization ON (Organization.id = AccountOrganization.organizationID) WHERE accountID = $1;", accountID); err != nil {
		return nil, err
	}
	return organizations, nil
}
