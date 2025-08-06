package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/isaacgr/portfolio/internal/api"
	"github.com/wneessen/go-mail"
)

type ContactDetails struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type MailClient struct {
	client *mail.Client
	user   string
	log    *slog.Logger
}

func NewMailClient(
	host string,
	user string,
	pass string,
	logger *slog.Logger,
) (*MailClient, error) {
	client, err := mail.NewClient(
		host,
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(
			mail.SMTPAuthAutoDiscover,
		),
		mail.WithUsername(user),
		mail.WithPassword(pass),
	)

	if err != nil {
		return nil, api.Error{
			Msg:  "Unable to send message",
			Code: http.StatusInternalServerError,
		}
	}
	return &MailClient{
		client: client,
		user:   user,
		log:    logger,
	}, nil
}

func NewContact(
	name string,
	email string,
	subject string,
	message string,
) ContactDetails {
	return ContactDetails{
		Name:    name,
		Email:   email,
		Subject: subject,
		Message: message,
	}
}

func (c *MailClient) SendEmail(
	contact ContactDetails,
) error {
	m := mail.NewMsg()
	m.Subject(fmt.Sprintf("<New contact message>: %s", contact.Subject))
	m.SetAddrHeader("From", c.user)
	m.SetAddrHeader("To", c.user)
	m.SetAddrHeader("Reply-To", contact.Email)

	messageBody := fmt.Sprintf(
		"Message from: %s\n\n%s",
		contact.Email,
		contact.Message,
	)
	m.SetBodyString(mail.TypeTextPlain, messageBody)

	if err := c.client.DialAndSend(m); err != nil {
		c.log.Error("Failed to send contact email.", "Error", err)
		return api.Error{
			Msg:  "Unable to send message",
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}
