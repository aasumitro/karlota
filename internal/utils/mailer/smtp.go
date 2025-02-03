package mailer

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/domodwyer/mailyak/v3"
	"karlota.aasumitro.id/internal/utils/ulid"
)

var _ Mailer = (*SMTPClient)(nil)

// SMTPClient defines a SMTP mail client structure that implements
// `mailer.Mailer` interface.
type SMTPClient struct {
	Host     string
	Port     int
	Username string
	Password string
	TLS      bool

	// LocalName is an optional domain name used for the EHLO/HELO exchange
	// (if not explicitly set, defaults to "localhost").
	//
	// This is required only by some SMTP servers, such as Gmail SMTP-relay.
	LocalName string
}

// Send implements `mailer.Mailer` interface.
func (c *SMTPClient) Send(m *Message) error {
	smtpAuth := smtp.PlainAuth("", c.Username, c.Password, c.Host)

	// create mail instance
	var yak *mailyak.MailYak
	if c.TLS {
		var tlsErr error
		yak, tlsErr = mailyak.NewWithTLS(fmt.Sprintf("%s:%d", c.Host, c.Port), smtpAuth, nil)
		if tlsErr != nil {
			return tlsErr
		}
	} else {
		yak = mailyak.New(fmt.Sprintf("%s:%d", c.Host, c.Port), smtpAuth)
	}

	if c.LocalName != "" {
		yak.LocalName(c.LocalName)
	}

	if m.From.Name != "" {
		yak.FromName(m.From.Name)
	}
	yak.From(m.From.Address)
	yak.Subject(m.Subject)
	yak.HTML().Set(m.HTML)

	if m.Text == "" {
		// try to generate a plain text version of the HTML
		if plain, err := HTML2Text(m.HTML); err == nil {
			yak.Plain().Set(plain)
		}
	} else {
		yak.Plain().Set(m.Text)
	}

	if len(m.To) > 0 {
		yak.To(addressesToStrings(m.To, true)...)
	}

	if len(m.Bcc) > 0 {
		yak.Bcc(addressesToStrings(m.Bcc, true)...)
	}

	if len(m.Cc) > 0 {
		yak.Cc(addressesToStrings(m.Cc, true)...)
	}

	// add attachments (if any)
	for name, data := range m.Attachments {
		yak.Attach(name, data)
	}

	// add custom headers (if any)
	var hasMessageID bool
	for k, v := range m.Headers {
		if strings.EqualFold(k, "Message-ID") {
			hasMessageID = true
		}
		yak.AddHeader(k, v)
	}
	if !hasMessageID {
		// add a default message id if missing
		fromParts := strings.Split(m.From.Address, "@")
		partLen := 2
		if len(fromParts) == partLen {
			yak.AddHeader("Message-ID", fmt.Sprintf("<%s@%s>",
				ulid.Generate("mail"),
				fromParts[1],
			))
		}
	}

	return yak.Send()
}
