package cmd

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/peragrin/api/auth"
	"gitlab.com/peragrin/api/communities"
	"gitlab.com/peragrin/api/db"
	"gitlab.com/peragrin/api/memberships"
	"gitlab.com/peragrin/api/organizations"
	"gitlab.com/peragrin/api/promotions"
	"gitlab.com/peragrin/api/service"
)

func serve() {
	log.SetFormatter(&log.JSONFormatter{})

	dbClient, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	storeClient, err := minio.New(viper.GetString("STORE_ENDPOINT"), viper.GetString("STORE_ACCESS_KEY"), viper.GetString("STORE_SECRET_KEY"), viper.GetBool("STORE_SECURE"))
	if err != nil {
		log.Fatal(err)
	}

	auth := auth.Init(dbClient, storeClient, viper.GetString("TOKEN_SECRET"), viper.GetString("LOCATIONIQ_API_KEY"), viper.GetString("APP_DOMAIN"), viper.GetString("MANDRILL_KEY"))
	organizations := organizations.Init(dbClient, storeClient)
	communities := communities.Init(dbClient, storeClient, viper.GetString("LOCATIONIQ_API_KEY"))
	memberships := memberships.Init(dbClient)
	promotions := promotions.Init(dbClient)

	r := mux.NewRouter()
	r.Handle("/auth/login", service.Handler(auth.LoginHandler))
	r.Handle("/auth/register", service.Handler(auth.RegisterHandler))
	r.Handle("/auth/forgot-password", service.Handler(auth.ForgotPasswordHandler))
	r.Handle("/auth/account", auth.RequiredMiddleware(auth.GetAccountHandler)).Methods(http.MethodGet)
	r.Handle("/auth/account", auth.RequiredMiddleware(auth.UpdateAccountHandler)).Methods(http.MethodPost)
	r.Handle("/auth/organizations", auth.RequiredMiddleware(auth.ListOrganizationsHandler)).Methods(http.MethodGet)
	r.Handle("/auth/organizations", auth.RequiredMiddleware(auth.CreateOrganizationHandler)).Methods(http.MethodPost)

	r.Handle("/communities", service.Handler(communities.ListHandler))
	r.Handle("/communities/{communityID:[0-9]+}", service.Handler(communities.GetHandler)).Methods(http.MethodGet)
	r.Handle("/communities/{communityID:[0-9]+}/organizations", service.Handler(communities.ListOrganizationsHandler)).Methods(http.MethodGet)
	r.Handle("/communities/{communityID:[0-9]+}/organizations", service.Handler(communities.CreateOrganizationHandler)).Methods(http.MethodPost)
	r.Handle("/communities/{communityID:[0-9]+}/posts", auth.RequiredMiddleware(communities.ListPostsHandler))
	r.Handle("/communities/{communityID:[0-9]+}/geo-json-overlays", service.Handler(communities.ListGeoJSONOverlaysHandler))
	r.Handle("/communities/{communityID:[0-9]+}/memberships", service.Handler(communities.ListMembershipsHandler)).Methods(http.MethodGet)
	r.Handle("/communities/{communityID:[0-9]+}/memberships", service.Handler(communities.CreateMembershipHandler)).Methods(http.MethodPost)

	r.Handle("/memberships/{membershipID:[0-9]+}/accounts", auth.RequiredMiddleware(memberships.ListAccountsHandler)).Methods(http.MethodGet)
	r.Handle("/memberships/{membershipID:[0-9]+}/accounts", auth.RequiredMiddleware(memberships.CreateAccountHandler)).Methods(http.MethodPost)

	r.Handle("/organizations/{organizationID:[0-9]+}", auth.RequiredMiddleware(organizations.UpdateHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}", auth.RequiredMiddleware(organizations.GetHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities", auth.RequiredMiddleware(organizations.ListCommunitiesHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities", auth.RequiredMiddleware(organizations.CreateCommunityHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities/{communityID:[0-9]+}", auth.RequiredMiddleware(organizations.JoinCommunityHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/posts", auth.RequiredMiddleware(organizations.CreatePostHandler))
	r.Handle("/organizations/{organizationID:[0-9]+}/hours", service.Handler(organizations.ListHoursHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/promotions", service.Handler(organizations.ListPromotionsHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/promotions", service.Handler(organizations.CreatePromotionHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/accounts", auth.RequiredMiddleware(organizations.ListAccountsHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/logo", auth.RequiredMiddleware(organizations.UploadLogoHandler)).Methods(http.MethodPost)

	r.Handle("/promotions/{promotionID:[0-9]+}/redeem", auth.RequiredMiddleware(promotions.RedeemHandler)).Methods(http.MethodPost)
	r.Handle("/promotions/{promotionID:[0-9]+}", auth.RequiredMiddleware(promotions.UpdateHandler)).Methods(http.MethodPut)
	r.Handle("/promotions/{promotionID:[0-9]+}", auth.RequiredMiddleware(promotions.DeleteHandler)).Methods(http.MethodDelete)

	log.Infof("initializing server: %s", viper.GetString("PORT"))

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%s", viper.GetString("PORT")),
		Handler:      http.TimeoutHandler(r, 15*time.Second, ""),
	}
	server.SetKeepAlivesEnabled(false)
	server.ListenAndServe()
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

	Serve.PersistentFlags().StringP("mandrill-key", "", "", "mandrill key")
	viper.BindPFlag("MANDRILL_KEY", Serve.PersistentFlags().Lookup("mandrill-key"))

	Serve.PersistentFlags().StringP("app-domain", "", "http://localhost:8080", "app domain")
	viper.BindPFlag("APP_DOMAIN", Serve.PersistentFlags().Lookup("app-domain"))
}
