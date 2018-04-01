package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jteppinette/peragrin-api/cmd"
)

const program = "peragrin-api"

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

	root.PersistentFlags().StringP("token-secret", "", "token-secret", "the secret used to sign the json web tokens")
	viper.BindPFlag("TOKEN_SECRET", root.PersistentFlags().Lookup("token-secret"))

	root.PersistentFlags().StringP("locationiq-api-key", "", "", "api key to access location iq api")
	viper.BindPFlag("LOCATIONIQ_API_KEY", root.PersistentFlags().Lookup("locationiq-api-key"))

	root.PersistentFlags().StringP("app-domain", "", "http://localhost:8080", "app domain")
	viper.BindPFlag("APP_DOMAIN", root.PersistentFlags().Lookup("app-domain"))

	root.PersistentFlags().StringP("mail-from", "", "notifications@peragrin.localhost", "mail from")
	viper.BindPFlag("MAIL_FROM", root.PersistentFlags().Lookup("mail-from"))

	root.PersistentFlags().StringP("mail-host", "", "0.0.0.0", "mail host")
	viper.BindPFlag("MAIL_HOST", root.PersistentFlags().Lookup("mail-host"))

	root.PersistentFlags().IntP("mail-port", "", 1025, "mail port")
	viper.BindPFlag("MAIL_PORT", root.PersistentFlags().Lookup("mail-port"))

	root.PersistentFlags().StringP("mail-password", "", "", "mail password")
	viper.BindPFlag("MAIL_PASSWORD", root.PersistentFlags().Lookup("mail-password"))

	root.PersistentFlags().StringP("mail-user", "", "", "mail user")
	viper.BindPFlag("MAIL_USER", root.PersistentFlags().Lookup("mail-user"))

	root.AddCommand(cmd.Migrate)
	root.AddCommand(cmd.Serve)
	root.AddCommand(cmd.AddSuperUser)
	root.AddCommand(cmd.SendTestMail)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
