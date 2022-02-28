package email

import (
	"net/mail"
)

type client struct {
	from string
	pwd  string
	host string
	port int
}

func NewClient(from, pwd, host string, port int) (*client, error) {
	_, err := mail.ParseAddress(from)
	if err != nil {
		return &client{}, err
	}

	return &client{
		from: from,
		pwd:  pwd,
		host: host,
		port: port,
	}, nil
}

func (c *client) Send(l Letter) error {
	if err := l.Validate(); err != nil {
		return err
	}

	return nil
}
