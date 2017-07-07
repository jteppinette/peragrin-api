package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/peragrin/api/db"
	"gitlab.com/peragrin/api/models"
)

func updateOrganizationLogos() {
	log.Info("initializing organization logo update")

	dbClient, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	storeClient, err := minio.New(viper.GetString("STORE_ENDPOINT"), viper.GetString("STORE_ACCESS_KEY"), viper.GetString("STORE_SECRET_KEY"), viper.GetBool("STORE_SECURE"))
	if err != nil {
		log.Fatal(err)
	}

	organizations, err := models.GetOrganizations(dbClient)
	if err != nil {
		log.Fatal(err)
	}

	for _, o := range organizations {
		if o.Logo == "" {
			continue
		}
		f, err := o.GetLogo(fmt.Sprintf("logos/%s", o.Logo), storeClient)
		if err != nil {
			log.Errorf("get logo: %v", err)
			continue
		}
		if err := o.UploadLogo(f, o.Logo, storeClient); err != nil {
			log.Errorf("upload logo: %v", err)
			continue
		}
		if err := o.RemoveLogo(storeClient); err != nil {
			log.Errorf("remove logo: %v", err)
			continue
		}
	}

	log.Info("completed successfully")
}

// UpdateOrganizationLogos is a cobra command that hooks into the db.migrate function.
var UpdateOrganizationLogos *cobra.Command

func init() {
	UpdateOrganizationLogos = &cobra.Command{
		Use: "update-organization-logos",
		Run: func(_ *cobra.Command, args []string) {
			updateOrganizationLogos()
		},
	}
}
