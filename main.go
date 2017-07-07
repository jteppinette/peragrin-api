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

		viper.AutomaticEnv()

		level, err := log.ParseLevel(viper.GetString("LOG_LEVEL"))
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(level)
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

	root.PersistentFlags().StringP("store-endpoint", "", "minio:9000", "store endpoint")
	viper.BindPFlag("STORE_ENDPOINT", root.PersistentFlags().Lookup("store-endpoint"))

	root.PersistentFlags().StringP("store-access-key", "", "access-key", "store access key")
	viper.BindPFlag("STORE_ACCESS_KEY", root.PersistentFlags().Lookup("store-access-key"))

	root.PersistentFlags().StringP("store-secret-key", "", "secret-key", "store secret key")
	viper.BindPFlag("STORE_SECRET_KEY", root.PersistentFlags().Lookup("store-secret-key"))

	root.PersistentFlags().BoolP("store-secure", "", false, "store secure")
	viper.BindPFlag("STORE_SECURE", root.PersistentFlags().Lookup("store-secure"))

	root.PersistentFlags().StringP("log-level", "l", "info", "log level [debug, info, warning, error, fatal, panic]")
	viper.BindPFlag("LOG_LEVEL", root.PersistentFlags().Lookup("log-level"))

	root.AddCommand(cmd.Migrate)
	root.AddCommand(cmd.Serve)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
