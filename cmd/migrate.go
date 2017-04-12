package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/peragrin/api/db"
)

func migrate() {
	client, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	log.Print("initializing migration")
	if err := db.Migrate(client); err != nil {
		log.Fatal(err)
	}
	log.Print("completed successfully")
}

var Migrate *cobra.Command

func init() {
	Migrate = &cobra.Command{
		Use: "migrate",
		Run: func(_ *cobra.Command, args []string) {
			migrate()
		},
	}
}
