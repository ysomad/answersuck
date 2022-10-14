package smtp

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	smail "github.com/xhit/go-simple-mail/v2"
)

const (
	defaultConnectTimeout = 10 * time.Second
	defaultSendTimeout    = 10 * time.Second
)

var (
	errEmptyHost     = errors.New("empty host")
	errEmptyFrom     = errors.New("empty from")
	errEmptyPassword = errors.New("empty password")
	errInvalidPort   = errors.New("invalid client port")
)

type client struct {
	srv  *smail.SMTPServer
	from string
}

type ClientOptions struct {
	Host           string
	Port           int
	From           string
	Password       string
	KeepAlive      bool
	ConnectTimeout time.Duration
	SendTimeout    time.Duration
}

func (opt *ClientOptions) validate() error {
	switch {
	case opt.Host == "":
		return errEmptyHost
	case opt.From == "":
		return errEmptyFrom
	case opt.Password == "":
		return errEmptyPassword
	case opt.Port == 0:
		return errInvalidPort
	case opt.ConnectTimeout == 0:
		opt.ConnectTimeout = defaultConnectTimeout
	case opt.SendTimeout == 0:
		opt.SendTimeout = defaultSendTimeout
	}

	_, err := mail.ParseAddress(opt.From)
	if err != nil {
		return err
	}

	return nil
}

func NewClient(opt *ClientOptions) (*client, error) {
	if err := opt.validate(); err != nil {
		return nil, fmt.Errorf("c.Validate(): %w", err)
	}

	c := smail.NewSMTPClient()

	c.Host = opt.Host
	c.Port = opt.Port
	c.Username = opt.From
	c.Password = opt.Password
	c.Encryption = smail.EncryptionSSLTLS

	c.KeepAlive = opt.KeepAlive
	c.ConnectTimeout = opt.ConnectTimeout
	c.SendTimeout = opt.SendTimeout

	return &client{
		srv:  c,
		from: opt.From,
	}, nil
}

func (c *client) SendEmail(ctx context.Context, to, subject, message string) error {
	msg := smail.NewMSG()

	msg.SetFrom(c.from).AddTo(to).SetSubject(subject)
	msg.SetBody(smail.TextHTML, message)

	cli, err := c.srv.Connect()
	if err != nil {
		return err
	}

	if err = msg.Send(cli); err != nil {
		return err
	}

	return nil
}
