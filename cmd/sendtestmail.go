package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jteppinette/peragrin-api/mail"
)

func sendTestMail(to []string) {
	client := mail.New(viper.GetString("MAIL_FROM"), viper.GetString("MAIL_HOST"), viper.GetInt("MAIL_PORT"), viper.GetString("MAIL_PASSWORD"), viper.GetString("MAIL_USER"))

	log.Info("sending mail")

	if err := client.Send(to, "test", "test"); err != nil {
		log.Fatal(err)
	}

	log.Info("mail sent")
}

// SendTestMail is a cobra command sends the recepients listed in the provided arguments a test email.
var SendTestMail *cobra.Command

func init() {
	SendTestMail = &cobra.Command{
		Use:  "sendtestmail",
		Args: cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			sendTestMail(args)
		},
	}
}
