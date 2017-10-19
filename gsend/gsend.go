package gsend

import (
	"bytes"
	"fmt"
	"net/smtp"
)

// Mail is a simple email.
type Mail struct {
	From    string
	To      string
	Subject string
	Message string
}

const server = "smtp.gmail.com"

// Send sends an email via gmail server.
func Send(m *Mail, password string) error {
	body := new(bytes.Buffer)
	fmt.Fprintf(body, "To: %s\r\n", m.To)
	fmt.Fprintf(body, "Subject: %s\r\n", m.Subject)
	fmt.Fprintf(body, "\r\n")
	fmt.Fprintf(body, "%s", m.Message)

	auth := smtp.PlainAuth("", m.From, password, server)
	return smtp.SendMail(
		fmt.Sprintf("%s:587", server),
		auth,
		m.From, []string{m.To},
		body.Bytes(),
	)
}
