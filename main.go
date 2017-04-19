package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/peragrin/api/cmd"
)

const program = "api"

var cfp string

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	cobra.OnInitialize(func() {
		if cfp != "" {
			viper.SetConfigFile(cfp)
			if err := viper.ReadInConfig(); err != nil {
				log.Fatal(err)
			}
			log.Infof("read configuration file: %s", cfp)
		}

		level, err := log.ParseLevel(viper.GetString("LOG_LEVEL"))
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(level)

		viper.AutomaticEnv()
	})

	root := &cobra.Command{
		Use:   program,
		Short: fmt.Sprintf("%s is a simple api that authentication endpoints and an authenticated resource", program),
	}

	root.PersistentFlags().StringVarP(&cfp, "config", "c", "", "config file path")

	root.PersistentFlags().StringP("db-host", "", "0.0.0.0", "db host")
	viper.BindPFlag("DB_HOST", root.PersistentFlags().Lookup("db-host"))

	root.PersistentFlags().StringP("db-user", "", "db", "db user")
	viper.BindPFlag("DB_USER", root.PersistentFlags().Lookup("db-user"))

	root.PersistentFlags().StringP("db-password", "", "secret", "db password")
	viper.BindPFlag("DB_PASSWORD", root.PersistentFlags().Lookup("db-password"))

	root.PersistentFlags().StringP("db-name", "", "db", "db name")
	viper.BindPFlag("DB_NAME", root.PersistentFlags().Lookup("db-name"))

	root.PersistentFlags().StringP("log-level", "l", "info", "log level [debug, info, warning, error, fatal, panic]")
	viper.BindPFlag("LOG_LEVEL", root.PersistentFlags().Lookup("log-level"))

	root.PersistentFlags().StringP("mapbox-api-key", "", "", "api key to access mapbox api")
	viper.BindPFlag("MAPBOX_API_KEY", root.PersistentFlags().Lookup("mapbox-api-key"))

	root.AddCommand(cmd.Migrate)
	root.AddCommand(cmd.CreateFixtureData)
	root.AddCommand(cmd.Serve)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
