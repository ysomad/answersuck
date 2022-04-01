package email

import (
	"net/mail"
	"time"

	smail "github.com/xhit/go-simple-mail/v2"
)

type client struct {
	server *smail.SMTPServer

	from string
}

func NewClient(from, pwd, host string, port int) (*client, error) {
	_, err := mail.ParseAddress(from)
	if err != nil {
		return &client{}, err
	}

	s := smail.NewSMTPClient()

	s.Host = host
	s.Port = port
	s.Username = from
	s.Password = pwd
	s.Encryption = smail.EncryptionSSLTLS

	s.KeepAlive = false
	s.ConnectTimeout = 10 * time.Second
	s.SendTimeout = 10 * time.Second

	return &client{
		server: s,
		from:   from,
	}, nil
}

func (c *client) Send(l Letter) error {
	if err := l.Validate(); err != nil {
		return err
	}

	msg := smail.NewMSG()
	msg.SetFrom(c.from).AddTo(l.To).SetSubject(l.Subject)
	msg.SetBody(smail.TextHTML, l.Message)

	cli, err := c.server.Connect()
	if err != nil {
		return err
	}

	if err = msg.Send(cli); err != nil {
		return err
	}

	return nil
}
