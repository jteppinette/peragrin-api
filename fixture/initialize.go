package fixture

import (
	"log"

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

func Initialize(client *sqlx.DB) error {

	log.Printf("\nCommunities\n------------------------------")
	if _, err := client.Exec("DELETE FROM communities"); err != nil {
		return err
	}
	for i, _ := range communities {
		community := &(communities[i])
		if err := community.Save(client); err != nil {
			return err
		}
		log.Printf("%+v", community)
	}
	log.Print()

	log.Printf("\nOrganizations & Users\n----------------------------")
	if _, err := client.Exec("DELETE FROM organizations"); err != nil {
		return err
	}
	if _, err := client.Exec("DELETE FROM users"); err != nil {
		return err
	}
	for i, _ := range organizations {
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

		log.Printf("%+v", organization)
		log.Printf("\t%+v", user)
		log.Print()
	}
	log.Print()

	return nil
}
