package fixture

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"gitlab.com/peragrin/api/models"
)

var users = []models.User{
	models.User{Username: "jteppinette"},
	models.User{Username: "sajohnson"},
}

var communities = []models.Community{
	models.Community{Name: "Midtown Atlanta"},
}

var organizations = []models.Organization{
	models.Organization{Name: "Midtown Chamber", Address: "50 Peachtree St", IsLeader: true},
	models.Organization{Name: "Papa John's Pizza", Address: "100 Peachtree Place", IsLeader: false},
}

// Initialize loads fixture data intot the current database. Any previously
// uploaded data will be deleted.
func Initialize(client *sqlx.DB) error {

	if _, err := client.Exec("DELETE FROM communities"); err != nil {
		return err
	}
	for i := range communities {
		community := &(communities[i])
		if err := community.Save(client); err != nil {
			return err
		}
		log.Infof("Created Community: %+v", community)
	}

	if _, err := client.Exec("DELETE FROM organizations"); err != nil {
		return err
	}
	if _, err := client.Exec("DELETE FROM users"); err != nil {
		return err
	}
	for i := range organizations {
		organization := &(organizations[i])
		organization.CommunityID = communities[0].ID
		if err := organization.Save(client); err != nil {
			return err
		}

		user := users[i]
		if err := user.SetPassword(user.Username); err != nil {
			return err
		}
		user.OrganizationID = organization.ID
		if err := user.Save(client); err != nil {
			return err
		}

		log.Infof("Created Organization: %+v", organization)
		log.Infof("Created User: %+v", user)
	}
	return nil
}
