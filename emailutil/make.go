package emailutil

import (
	"bytes"
	"fmt"
	"time"
)

// Address is an email address.
type Address struct {
	Name  string
	Email string
}

// NewAddress creates a new email address.
func NewAddress(name, email string) *Address {
	return &Address{Name: name, Email: email}
}

func (a *Address) String() string {
	if a.Name == "" {
		return a.Email
	}
	return fmt.Sprintf("%s <%s>", a.Name, a.Email)
}

// Header is an email header.
type Header struct {
	From    *Address
	To      *Address
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
	printHeader(b, "From", h.From.String())
	printHeader(b, "To", h.To.String())
	printHeader(b, "Subject", h.Subject)
	printHeader(b, "MIME-Version", "1.0;")
	printHeader(b, "Content-Type", `test/html; charset="UTF-8"`)
	fmt.Fprint(b, "\r\n\r\n")
	b.Write(body)
	return b.Bytes()
}
