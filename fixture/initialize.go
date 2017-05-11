package fixture

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"gitlab.com/peragrin/api/models"
)

var accounts = []models.Account{
	models.Account{Email: "jteppinette@jteppinette.com"},
	models.Account{Email: "sajohnson@sajohnson.com"},
}

var communities = []models.Community{
	models.Community{Name: "Midtown Atlanta"},
}

var organizations = []models.Organization{
	models.Organization{Name: "Metro Atlanta Chamber", Address: "191 Peachtree Tower, 191 Peachtree St NE #3400, Atlanta, GA 30303", Longitude: -84.38642, Latitude: 33.759115, Leader: true, Enabled: true},
	models.Organization{Name: "Bobby Dodd Stadium", Address: "North Avenue NW, Atlanta, GA 30313", Longitude: -84.3903448, Latitude: 33.7712937, Leader: false},
}

// Initialize loads fixture data intot the current database. Any previously
// uploaded data will be deleted.
func Initialize(client *sqlx.DB) error {

	if _, err := client.Exec("DELETE FROM Community"); err != nil {
		return err
	}
	for i := range communities {
		community := &(communities[i])
		if err := community.Save(client); err != nil {
			return err
		}
		log.Infof("Created Community: %+v", community)
	}

	if _, err := client.Exec("DELETE FROM Organization"); err != nil {
		return err
	}
	if _, err := client.Exec("DELETE FROM Account"); err != nil {
		return err
	}
	for i := range organizations {
		organization := &(organizations[i])
		organization.CommunityID = communities[0].ID
		if err := organization.Save(client); err != nil {
			return err
		}

		account := accounts[i]
		if err := account.SetPassword(strings.Split(account.Email, "@")[0]); err != nil {
			return err
		}
		if err := account.Save(client); err != nil {
			return err
		}

		if err := account.AddOperator(organization.ID, client); err != nil {
			return err
		}

		post := models.Post{OrganizationID: organization.ID, Content: "We just got setup with Peragrin!"}
		if err := post.Save(client); err != nil {
			return err
		}

		log.Infof("Created Organization: %+v", organization)
		log.Infof("Created Post: %+v", post)
		log.Infof("Created Account: %+v", account)
	}
	return nil
}
