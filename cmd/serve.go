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
)

func serve() {
	log.SetFormatter(&log.JSONFormatter{})

	client, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	auth := auth.Init(client, viper.GetString("TOKEN_SECRET"))
	organizations := organizations.Init(client, viper.GetString("LOCATIONIQ_API_KEY"))
	communities := communities.Init(client)

	r := mux.NewRouter()
	r.Handle("/auth/login", service.Handler(auth.LoginHandler))
	r.Handle("/auth/register", service.Handler(auth.RegisterHandler))
	r.Handle("/auth/account", auth.RequiredMiddleware(auth.AccountHandler))
	r.Handle("/auth/organizations", auth.RequiredMiddleware(auth.OrganizationsHandler))

	r.Handle("/communities", service.Handler(communities.ListHandler))
	r.Handle("/communities/{communityID:[0-9]+}/organizations", auth.RequiredMiddleware(communities.ListOrganizationsHandler))
	r.Handle("/communities/{communityID:[0-9]+}/posts", auth.RequiredMiddleware(communities.ListPostsHandler))

	r.Handle("/organizations", auth.RequiredMiddleware(organizations.ListHandler)).Methods(http.MethodGet)
	r.Handle("/organizations", auth.RequiredMiddleware(organizations.CreateHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}", auth.RequiredMiddleware(organizations.GetHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}", auth.RequiredMiddleware(organizations.UpdateHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities", auth.RequiredMiddleware(organizations.ListCommunitiesHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities", auth.RequiredMiddleware(organizations.BulkJoinCommunityHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/posts", auth.RequiredMiddleware(organizations.CreatePostHandler))

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

	Serve.PersistentFlags().StringP("locationiq-api-key", "", "", "api key to access location iq api")
	viper.BindPFlag("LOCATIONIQ_API_KEY", Serve.PersistentFlags().Lookup("locationiq-api-key"))
}
