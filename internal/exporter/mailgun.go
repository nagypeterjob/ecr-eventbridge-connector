package exporter

import (
	"fmt"
	"strings"

	"github.com/mailgun/mailgun-go"
	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/formatter"
)

// MailgunExporter lets you send reports via email
type MailgunExporter struct {
	client     *mailgun.MailgunImpl
	from       string
	name       string
	recipients string
}

// NewMailgunExporter .
func NewMailgunExporter(name string, recipients string, from string, apiKey string) *MailgunExporter {
	domain := domain(from)

	return &MailgunExporter{
		client:     mailgun.NewMailgun(domain, apiKey),
		from:       from,
		recipients: recipients,
		name:       name,
	}
}

// Name .
func (m MailgunExporter) Name() string {
	return m.name
}

// Send sends report via Mailgun
func (m MailgunExporter) Send(msg *formatter.Message) error {

	mail := m.client.NewMessage(
		m.from,
		fmt.Sprintf("New vulnerabilites found in %s", msg.RepositoryName),
		msg.Title+msg.Body+msg.Link,
	)

	recipientList := strings.Split(m.recipients, ",")

	for _, user := range recipientList {
		if err := mail.AddRecipient(user); err != nil {
			return err
		}
	}

	if _, _, err := m.client.Send(mail); err != nil {
		return err
	}

	return nil

}

// extracts domain from email address
func domain(email string) string {
	if email == "" {
		return ""
	}

	if strings.Count(email, "@") > 1 {
		return ""
	}

	atPosition := strings.Index(email, "@")
	return email[atPosition+1:]
}
