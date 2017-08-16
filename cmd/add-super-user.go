package cmd

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/mattbaird/gochimp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/peragrin/api/db"
	"gitlab.com/peragrin/api/models"
)

func addSuperUser() {
	log.SetFormatter(&log.JSONFormatter{})

	dbClient, err := db.Client(viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	mailClient, err := gochimp.NewMandrill(viper.GetString("MANDRILL_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter super user email address: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	account := models.Account{Email: scanner.Text()}

	if existing, err := models.GetAccountByEmail(account.Email, dbClient); err != nil {
		log.Fatal(err)
	} else if existing != nil {
		log.WithFields(log.Fields{"email": existing.Email, "accountID": existing.ID}).Info("account found - updating to super user status")
		if err := existing.SetIsSuper(true, dbClient); err != nil {
			log.WithFields(log.Fields{"email": existing.Email, "error": err.Error()}).Fatal(errors.New("account set is super"))
		}
		log.WithFields(log.Fields{
			"email":   existing.Email,
			"isSuper": existing.IsSuper,
		}).Info("updated account to super user status")
		return
	}

	log.WithFields(log.Fields{"email": account.Email}).Info("account not found - creating new super user")

	if err := account.Save(dbClient); err != nil {
		log.WithFields(log.Fields{"email": account.Email, "error": err.Error()}).Fatal(errors.New("account creation"))
	}

	if err := account.SetIsSuper(true, dbClient); err != nil {
		log.WithFields(log.Fields{"email": account.Email, "error": err.Error()}).Fatal(errors.New("update super user status"))
	}

	if err := account.SendActivationEmail("/communities", viper.GetString("APP_DOMAIN"), viper.GetString("TOKEN_SECRET"), "Super User", mailClient); err != nil {
		log.WithFields(log.Fields{"email": account.Email, "error": err.Error()}).Fatal(errors.New("account activation email"))
	}

	log.WithFields(log.Fields{
		"email":   account.Email,
		"isSuper": account.IsSuper,
	}).Info("account activation email sent")
}

// AddSuperUser adds a new super user. If the user already exists they will have the isSuper flag
// set to true. Otherwise, a new user will be created with isSuper: true, and an activation email
// will be sent.
var AddSuperUser *cobra.Command

func init() {
	AddSuperUser = &cobra.Command{
		Use: "add-super-user",
		Run: func(_ *cobra.Command, args []string) {
			addSuperUser()
		},
	}
}
