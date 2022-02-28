package email

import (
	"bytes"
	"errors"
	"html/template"
	"net/mail"
)

var (
	ErrEmailEmptyTo      = errors.New("to cannot be empty")
	ErrEmailEmptySubject = errors.New("subject cannot be empty")
	ErrEmailEmptyBody    = errors.New("body cannot be empty")
)

type Letter struct {
	To      string
	Subject string
	body    string
}

func (l *Letter) SetBodyFromTemplate(filename string, data interface{}) error {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	l.body = buf.String()

	return nil
}

func (l *Letter) Validate() error {
	if l.To == "" {
		return ErrEmailEmptyTo
	}

	_, err := mail.ParseAddress(l.To)
	if err != nil {
		return err
	}

	if l.Subject == "" {
		return ErrEmailEmptySubject
	}

	if l.body == "" {
		return ErrEmailEmptyBody
	}

	return nil
}
