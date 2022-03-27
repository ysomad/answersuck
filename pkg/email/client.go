package email

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

type client struct {
	from string
	host string
	port int
	auth smtp.Auth
}

func NewClient(from, pwd, host string, port int) (*client, error) {
	_, err := mail.ParseAddress(from)
	if err != nil {
		return &client{}, err
	}

	a := smtp.PlainAuth("", from, pwd, host)

	return &client{
		from: from,
		auth: a,
		host: host,
		port: port,
	}, nil
}

func (c *client) Send(l Letter) error {
	if err := l.Validate(); err != nil {
		return err
	}

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%d", c.host, c.port),
		c.auth,
		c.from,
		[]string{l.To},
		[]byte(l.Message),
	); err != nil {
		return err
	}

	return nil
}
