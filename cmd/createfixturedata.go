package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/peragrin/api/db"
	"gitlab.com/peragrin/api/fixture"
)

func createfixturedata() {
	client, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("creating fixture data - this will remove all data from the database")
	if err := fixture.Initialize(client); err != nil {
		log.Fatal(err)
	}
	log.Info("completed successfully")
}

// CreateFixtureData is a cobra command that hooks into the fixture.Initialize function.
var CreateFixtureData *cobra.Command

func init() {
	CreateFixtureData = &cobra.Command{
		Use: "createfixturedata",
		Run: func(_ *cobra.Command, args []string) {
			createfixturedata()
		},
	}
}
