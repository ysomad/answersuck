package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
)

type Email struct {
	To      string
	Subject string
	message string
}

func (e *Email) Message() string { return e.message }

func (e *Email) setMessageFromTemplate(filename string, data interface{}) error {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, &data); err != nil {
		return err
	}

	e.message = buf.String()

	return nil
}

func (e *Email) validate() error {
	if e.To == "" {
		return ErrEmptyTo
	}

	_, err := mail.ParseAddress(e.To)
	if err != nil {
		return fmt.Errorf("mail.ParseAddress: %w", err)
	}

	if e.Subject == "" {
		return ErrEmptySubject
	}

	if e.message == "" {
		return ErrEmptyMessage

	}

	return nil
}
