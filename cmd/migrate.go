package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	mattes "github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jteppinette/peragrin-api/db"

	// This import register the file source driver for mattes/migrate.
	_ "github.com/mattes/migrate/source/file"
)

func migrate() {
	client, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("initializing migration")

	driver, err := postgres.WithInstance(client.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "create driver"))
	}
	m, err := mattes.NewWithDatabaseInstance(fmt.Sprintf("file://%s", viper.GetString("MIGRATIONS_DIRECTORY")), "postgres", driver)
	if err != nil {
		log.Fatal(errors.Wrap(err, "initialize migration instance"))
	}

	if err := m.Up(); err != nil {
		log.Fatal(errors.Wrap(err, "migrate"))
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

	Migrate.PersistentFlags().StringP("migrations-directory", "m", "", "absolute path to migrations directory")
	viper.BindPFlag("MIGRATIONS_DIRECTORY", Migrate.PersistentFlags().Lookup("migrations-directory"))
}
