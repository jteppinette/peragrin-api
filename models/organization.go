package models

import (
	"github.com/jmoiron/sqlx"
)

type Organizations []Organization

type Organization struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	IsLeader    bool   `json:"isLeader"`
	CommunityID int    `json:"communityID"`
}

func (o *Organization) Save(client *sqlx.DB) error {
	if o.ID != 0 {
		return client.Get(o, "UPDATE organizations SET name = $2, address = $3, isLeader = $4, communityID = $4 WHERE id = $1 RETURNING *;", o.ID, o.Name, o.Address, o.IsLeader, o.CommunityID)
	} else {
		return client.Get(o, "INSERT INTO organizations (name, address, isLeader, communityID) VALUES ($1, $2, $3, $4) RETURNING *;", o.Name, o.Address, o.IsLeader, o.CommunityID)
	}
}

func ListOrganizations(client *sqlx.DB) (Organizations, error) {
	organizations := Organizations{}
	if err := client.Select(&organizations, "SELECT * FROM organizations;"); err != nil {
		return nil, err
	}
	return organizations, nil
}

func GetOrganizationByID(id int, client *sqlx.DB) (*Organization, error) {
	o := &Organization{}
	if err := client.Get(o, "SELECT * FROM organizations WHERE id = $1;", id); err != nil {
		return nil, err
	}
	return o, nil
}
