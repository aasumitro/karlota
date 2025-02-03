package mailer

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/mail"
	"os/exec"
	"strings"
)

// Message defines a generic email message struct.
type Message struct {
	From        mail.Address         `json:"from"`
	To          []mail.Address       `json:"to"`
	Bcc         []mail.Address       `json:"bcc"`
	Cc          []mail.Address       `json:"cc"`
	Subject     string               `json:"subject"`
	HTML        string               `json:"html"`
	Text        string               `json:"text"`
	Headers     map[string]string    `json:"headers"`
	Attachments map[string]io.Reader `json:"attachments"`
}

// Mailer defines a base mail client interface.
type Mailer interface {
	// Send sends an email with the provided Message.
	Send(message *Message) error
}

// addressesToStrings converts the provided address to a list of serialized RFC 5322 strings.
//
// To export only the email part of mail.Address, you can set withName to false.
func addressesToStrings(addresses []mail.Address, withName bool) []string {
	result := make([]string, len(addresses))
	for i, addr := range addresses {
		if withName && addr.Name != "" {
			result[i] = addr.String()
		} else {
			// keep only the email part to avoid wrapping in angle-brackets
			result[i] = addr.Address
		}
	}

	return result
}

// Sendmail implements [mailer.Mailer] interface and defines a mail
// client that sends emails via the "sendmail" *nix command.
//
// This client is usually recommended only for development and testing.
type Sendmail struct{}

// Send implements `mailer.Mailer` interface.
func (c *Sendmail) Send(m *Message) error {
	toAddresses := addressesToStrings(m.To, false)
	headers := make(http.Header)
	headers.Set("Subject", mime.QEncoding.Encode("utf-8", m.Subject))
	headers.Set("From", m.From.String())
	headers.Set("Content-Type", "text/html; charset=UTF-8")
	headers.Set("To", strings.Join(toAddresses, ","))

	cmdPath, err := findSendmailPath()
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	// write
	// ---
	if err := headers.Write(&buffer); err != nil {
		return err
	}
	if _, err := buffer.Write([]byte("\r\n")); err != nil {
		return err
	}
	if m.HTML != "" {
		if _, err := buffer.Write([]byte(m.HTML)); err != nil {
			return err
		}
	} else {
		if _, err := buffer.Write([]byte(m.Text)); err != nil {
			return err
		}
	}
	// ---

	sendmail := exec.Command(cmdPath, strings.Join(toAddresses, ","))
	sendmail.Stdin = &buffer
	return sendmail.Run()
}

func findSendmailPath() (string, error) {
	options := []string{
		"/usr/sbin/sendmail",
		"/usr/bin/sendmail",
		"sendmail",
	}

	for _, option := range options {
		path, err := exec.LookPath(option)
		if err == nil {
			return path, nil
		}
	}

	return "", errors.New("failed to locate a sendmail executable path")
}
