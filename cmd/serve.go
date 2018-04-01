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

	"github.com/jteppinette/peragrin-api/accounts"
	"github.com/jteppinette/peragrin-api/auth"
	"github.com/jteppinette/peragrin-api/communities"
	"github.com/jteppinette/peragrin-api/db"
	"github.com/jteppinette/peragrin-api/geo"
	"github.com/jteppinette/peragrin-api/mail"
	"github.com/jteppinette/peragrin-api/memberships"
	"github.com/jteppinette/peragrin-api/organizations"
	"github.com/jteppinette/peragrin-api/promotions"
	"github.com/jteppinette/peragrin-api/service"
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

	mailClient := mail.New(viper.GetString("MAIL_FROM"), viper.GetString("MAIL_HOST"), viper.GetInt("MAIL_PORT"), viper.GetString("MAIL_PASSWORD"), viper.GetString("MAIL_USER"))

	auth := auth.Init(dbClient, mailClient, viper.GetString("TOKEN_SECRET"), viper.GetString("APP_DOMAIN"))
	accounts := accounts.Init(dbClient, storeClient, mailClient, viper.GetString("TOKEN_SECRET"), viper.GetString("APP_DOMAIN"))
	organizations := organizations.Init(dbClient, storeClient, mailClient, viper.GetString("TOKEN_SECRET"), viper.GetString("APP_DOMAIN"))
	geo := geo.Init(viper.GetString("LOCATIONIQ_API_KEY"))
	communities := communities.Init(dbClient, storeClient, mailClient, viper.GetString("TOKEN_SECRET"), viper.GetString("APP_DOMAIN"))
	memberships := memberships.Init(dbClient, mailClient, viper.GetString("TOKEN_SECRET"), viper.GetString("APP_DOMAIN"))
	promotions := promotions.Init(dbClient)

	r := mux.NewRouter()
	r.Handle("/auth/login", service.Handler(auth.LoginHandler)).Methods(http.MethodPost)
	r.Handle("/auth/register", service.Handler(auth.RegisterHandler)).Methods(http.MethodPost)
	r.Handle("/auth/forgot-password", service.Handler(auth.ForgotPasswordHandler)).Methods(http.MethodPost)
	r.Handle("/auth/set-password", auth.RequiredMiddleware(auth.SetPasswordHandler)).Methods(http.MethodPost)
	r.Handle("/auth/activate", auth.RequiredMiddleware(auth.ActivateHandler)).Methods(http.MethodPost)

	r.Handle("/geo", auth.RequiredMiddleware(geo.LookupHandler)).Methods(http.MethodPost)

	r.Handle("/accounts", auth.RequiredMiddleware(accounts.ListHandler)).Methods(http.MethodGet)
	r.Handle("/accounts/{accountID:[0-9]+}", auth.RequiredMiddleware(accounts.UpdateAccountHandler)).Methods(http.MethodPut)
	r.Handle("/accounts/{accountID:[0-9]+}/forgot-password", auth.RequiredMiddleware(accounts.ForgotPasswordHandler)).Methods(http.MethodPost)
	r.Handle("/accounts/{accountID:[0-9]+}/organizations", auth.RequiredMiddleware(accounts.ListOrganizationsHandler)).Methods(http.MethodGet)
	r.Handle("/accounts/{accountID:[0-9]+}/organizations", auth.RequiredMiddleware(accounts.CreateOrganizationHandler)).Methods(http.MethodPost)
	r.Handle("/accounts/{accountID:[0-9]+}/communities", auth.RequiredMiddleware(accounts.ListCommunitiesHandler)).Methods(http.MethodGet)
	r.Handle("/accounts/{accountID:[0-9]+}/promotions/{promotionID:[0-9]+}", auth.RequiredMiddleware(accounts.ListPromotionRedemptionsHandler)).Methods(http.MethodGet)
	r.Handle("/accounts/{accountID:[0-9]+}/promotions", auth.RequiredMiddleware(accounts.ListRedemptionsHandler)).Methods(http.MethodGet)
	r.Handle("/accounts/{accountID:[0-9]+}/memberships", auth.RequiredMiddleware(accounts.ListMembershipsByCommunityHandler)).Methods(http.MethodGet)

	r.Handle("/communities", service.Handler(communities.ListHandler)).Methods(http.MethodGet)
	r.Handle("/communities", auth.RequiredMiddleware(communities.CreateHandler)).Methods(http.MethodPost)
	r.Handle("/communities/{communityID:[0-9]+}", service.Handler(communities.GetHandler)).Methods(http.MethodGet)
	r.Handle("/communities/{communityID:[0-9]+}", auth.RequiredMiddleware(communities.DeleteHandler)).Methods(http.MethodDelete)
	r.Handle("/communities/{communityID:[0-9]+}", auth.RequiredMiddleware(communities.UpdateHandler)).Methods(http.MethodPut)
	r.Handle("/communities/{communityID:[0-9]+}/organizations", service.Handler(communities.ListOrganizationsHandler)).Methods(http.MethodGet)
	r.Handle("/communities/{communityID:[0-9]+}/organizations", auth.RequiredMiddleware(communities.CreateOrganizationHandler)).Methods(http.MethodPost)
	r.Handle("/communities/{communityID:[0-9]+}/posts", auth.RequiredMiddleware(communities.ListPostsHandler))
	r.Handle("/communities/{communityID:[0-9]+}/geo-json-overlays", service.Handler(communities.ListGeoJSONOverlaysHandler))
	r.Handle("/communities/{communityID:[0-9]+}/memberships", auth.RequiredMiddleware(communities.ListMembershipsHandler)).Methods(http.MethodGet)
	r.Handle("/communities/{communityID:[0-9]+}/memberships", auth.RequiredMiddleware(communities.CreateMembershipHandler)).Methods(http.MethodPost)
	r.Handle("/communities/{communityID:[0-9]+}/accounts", auth.RequiredMiddleware(communities.BulkAddAccountsHandler)).Methods(http.MethodPost).Headers("X-Action", "bulk")

	r.Handle("/memberships/{membershipID:[0-9]+}", auth.RequiredMiddleware(memberships.GetHandler)).Methods(http.MethodGet)
	r.Handle("/memberships/{membershipID:[0-9]+}", auth.RequiredMiddleware(memberships.UpdateHandler)).Methods(http.MethodPut)
	r.Handle("/memberships/{membershipID:[0-9]+}", auth.RequiredMiddleware(memberships.DeleteHandler)).Methods(http.MethodDelete)
	r.Handle("/memberships/{membershipID:[0-9]+}/accounts", auth.RequiredMiddleware(memberships.ListAccountsHandler)).Methods(http.MethodGet)
	r.Handle("/memberships/{membershipID:[0-9]+}/accounts", auth.RequiredMiddleware(memberships.BulkAddAccountsHandler)).Methods(http.MethodPost).Headers("X-Action", "bulk")
	r.Handle("/memberships/{membershipID:[0-9]+}/accounts", auth.RequiredMiddleware(memberships.AddAccountHandler)).Methods(http.MethodPost)
	r.Handle("/memberships/{membershipID:[0-9]+}/accounts/{accountID:[0-9]+}", auth.RequiredMiddleware(memberships.RemoveAccountHandler)).Methods(http.MethodDelete)
	r.Handle("/memberships/{membershipID:[0-9]+}/accounts/{accountID:[0-9]+}", auth.RequiredMiddleware(memberships.UpdateAccountHandler)).Methods(http.MethodPut)

	r.Handle("/organizations/{organizationID:[0-9]+}", auth.RequiredMiddleware(organizations.UpdateHandler)).Methods(http.MethodPut)
	r.Handle("/organizations/{organizationID:[0-9]+}", auth.RequiredMiddleware(organizations.GetHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities", auth.RequiredMiddleware(organizations.ListCommunitiesHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities", auth.RequiredMiddleware(organizations.CreateCommunityHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities/{communityID:[0-9]+}", auth.RequiredMiddleware(organizations.JoinCommunityHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/communities/{communityID:[0-9]+}", auth.RequiredMiddleware(organizations.RemoveCommunityHandler)).Methods(http.MethodDelete)
	r.Handle("/organizations/{organizationID:[0-9]+}/posts", auth.RequiredMiddleware(organizations.CreatePostHandler))
	r.Handle("/organizations/{organizationID:[0-9]+}/hours", service.Handler(organizations.ListHoursHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/promotions", service.Handler(organizations.ListPromotionsHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/promotions", service.Handler(organizations.CreatePromotionHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/accounts", auth.RequiredMiddleware(organizations.ListAccountsHandler)).Methods(http.MethodGet)
	r.Handle("/organizations/{organizationID:[0-9]+}/accounts", auth.RequiredMiddleware(organizations.AddAccountHandler)).Methods(http.MethodPost)
	r.Handle("/organizations/{organizationID:[0-9]+}/accounts/{accountID:[0-9]+}", auth.RequiredMiddleware(organizations.RemoveAccountHandler)).Methods(http.MethodDelete)
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
}
