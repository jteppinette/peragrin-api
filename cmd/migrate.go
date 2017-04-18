package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/peragrin/api/db"
)

func migrate() {
	client, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("initializing migration")
	if err := db.Migrate(client); err != nil {
		log.Fatal(err)
	}
	log.Info("completed successfully")
}

// Migrate is a cobra command that hooks into the db.migrate function.
var Migrate *cobra.Command

func init() {
	Migrate = &cobra.Command{
		Use: "migrate",
		Run: func(_ *cobra.Command, args []string) {
			migrate()
		},
	}
}
