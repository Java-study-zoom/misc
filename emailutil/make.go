package emailutil

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// Address creates a new email address.
func Address(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

// Header is an email header.
type Header struct {
	From    string
	To      string
	Subject string
	Time    time.Time
}

func printHeader(b *bytes.Buffer, k, v string) {
	fmt.Fprintf(b, "%s: %s\r\n", k, v)
}

// Make creates an email with the given header and body.
func Make(h *Header, body []byte) []byte {
	b := new(bytes.Buffer)
	printHeader(b, "Date", h.Time.String())
	printHeader(b, "From", h.From)
	printHeader(b, "To", h.To)
	printHeader(b, "Subject", h.Subject)
	printHeader(b, "MIME-Version", "1.0;")
	printHeader(b, "Content-Type", `text/html; charset="UTF-8"`)
	fmt.Fprint(b, "\r\n")
	b.Write(body)
	return b.Bytes()
}

// TemplateMake creates an email using the given template
func TemplateMake(
	h *Header, t *template.Template, dat interface{},
) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, dat); err != nil {
		return nil, err
	}
	return Make(h, buf.Bytes()), nil
}
