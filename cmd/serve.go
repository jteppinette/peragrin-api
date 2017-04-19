package cmd

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/peragrin/api/auth"
	"gitlab.com/peragrin/api/communities"
	"gitlab.com/peragrin/api/db"
	"gitlab.com/peragrin/api/organizations"
	"gitlab.com/peragrin/api/service"
	"gitlab.com/peragrin/api/users"
)

func serve() {
	log.SetFormatter(&log.JSONFormatter{})

	client, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	auth := auth.Init(client, viper.GetString("TOKEN_SECRET"))
	users := users.Init(client)
	organizations := organizations.Init(client)
	communities := communities.Init(client)

	r := mux.NewRouter()
	r.Handle("/login", service.Handler(auth.LoginHandler))
	r.Handle("/register", service.Handler(auth.RegisterHandler))
	r.Handle("/user", auth.RequiredMiddleware(auth.UserHandler))
	r.Handle("/communities", service.Handler(communities.ListHandler))

	r.Handle("/communities/{communityID:[0-9]+}/organizations", auth.RequiredMiddleware(communities.ListOrganizationsHandler))

	r.Handle("/users", auth.RequiredMiddleware(users.ListHandler))

	r.Handle("/organizations", auth.RequiredMiddleware(organizations.ListHandler))
	r.Handle("/organizations/{organizationID:[0-9]+}", auth.RequiredMiddleware(organizations.GetHandler))

	log.Infof("initializing server: %s", viper.GetString("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("PORT")), r)
}

// Serve instantiates the API server.
var Serve *cobra.Command

func init() {
	Serve = &cobra.Command{
		Use: "serve",
		Run: func(_ *cobra.Command, args []string) {
			serve()
		},
	}

	Serve.PersistentFlags().StringP("port", "", "8000", "port that the api will listen on")
	viper.BindPFlag("PORT", Serve.PersistentFlags().Lookup("port"))

	Serve.PersistentFlags().StringP("token-secret", "", "token-secret", "the secret used to sign the json web tokens")
	viper.BindPFlag("TOKEN_SECRET", Serve.PersistentFlags().Lookup("token-secret"))
}
