package cmd

import (
	log "github.com/Sirupsen/logrus"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/peragrin/api/db"
	"gitlab.com/peragrin/api/fixture"
)

func createfixturedata() {
	dbClient, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	storeClient, err := minio.New(viper.GetString("STORE_ENDPOINT"), viper.GetString("STORE_ACCESS_KEY"), viper.GetString("STORE_SECRET_KEY"), viper.GetBool("STORE_SECURE"))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("creating fixture data - this will remove all data from the database")
	if err := fixture.Initialize(dbClient, storeClient, viper.GetString("DIR")); err != nil {
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

	CreateFixtureData.PersistentFlags().StringP("dir", "", "/etc/peragrin/fixture", "absolure directory of fixture files")
	viper.BindPFlag("DIR", CreateFixtureData.PersistentFlags().Lookup("dir"))
}
