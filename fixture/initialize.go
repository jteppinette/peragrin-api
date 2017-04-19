package fixture

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codingsince1985/geo-golang/mapbox"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gitlab.com/peragrin/api/models"
)

var users = []models.User{
	models.User{Email: "jteppinette@jteppinette.com"},
	models.User{Email: "sajohnson@sajohnson.com"},
}

var communities = []models.Community{
	models.Community{Name: "Midtown Atlanta"},
}

var organizations = []models.Organization{
	models.Organization{Name: "Publix at The Plaza Midtown", Address: "950 W Peachtree St NE, Atlanta, GA 30309", IsLeader: true},
	models.Organization{Name: "Bobby Dodd Stadium", Address: "North Avenue NW, Atlanta, GA 30313", IsLeader: false},
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
		if err := organization.SetGeo(mapbox.Geocoder(viper.GetString("MAPBOX_API_KEY"))); err != nil {
			return err
		}
		if err := organization.Save(client); err != nil {
			return err
		}

		user := users[i]
		if err := user.SetPassword(strings.Split(user.Email, "@")[0]); err != nil {
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
