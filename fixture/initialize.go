package fixture

import (
	"log"

	"github.com/jmoiron/sqlx"
	"gitlab.com/peragrin/api/models"
)

var users = []models.User{
	models.User{Username: "jteppinette"},
	models.User{Username: "sajohnson"},
	models.User{Username: "mbfulton"},
}

func Initialize(client *sqlx.DB) error {
	if _, err := client.Exec("DELETE FROM users"); err != nil {
		return err
	}
	log.Printf("\nUsers\n------------------------------------")
	for _, user := range users {
		if err := user.SetPassword(user.Username); err != nil {
			return err
		}
		if err := user.Save(client); err != nil {
			return err
		}
		log.Printf("%+v", user)
	}
	log.Print()
	return nil
}
