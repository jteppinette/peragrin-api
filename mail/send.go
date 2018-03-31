package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

func (c *Config) Send(to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", c.User, c.Password, c.Host)
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", c.Host, c.Port),
		auth,
		c.From,
		to,
		[]byte(fmt.Sprintf(
			"To: %s\r\nSubject: %s\r\n\r\n%s\r\n",
			strings.Join(to, ", "),
			subject,
			body,
		)),
	)
}
