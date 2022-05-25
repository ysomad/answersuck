package smtp

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/answersuck/vault/internal/app/email"

	smail "github.com/xhit/go-simple-mail/v2"
)

const (
	defaultConnectTimeout = 10 * time.Second
	defaultSendTimeout    = 10 * time.Second
)

var (
	ErrEmptyHost     = errors.New("empty host")
	ErrEmptyFrom     = errors.New("empty from")
	ErrEmptyPassword = errors.New("empty password")
	ErrInvalidPort   = errors.New("invalid client port")
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

func (opt *ClientOptions) Validate() error {
	switch {
	case opt.Host == "":
		return ErrEmptyHost
	case opt.From == "":
		return ErrEmptyFrom
	case opt.Password == "":
		return ErrEmptyPassword
	case opt.Port == 0:
		return ErrInvalidPort
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
	if err := opt.Validate(); err != nil {
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

func (c *client) SendEmail(ctx context.Context, e email.Email) error {
	if err := e.Validate(); err != nil {
		return fmt.Errorf("e.Validate: %w", err)
	}

	msg := smail.NewMSG()

	msg.SetFrom(c.from).AddTo(e.To).SetSubject(e.Subject)
	msg.SetBody(smail.TextHTML, e.Message())

	cli, err := c.srv.Connect()
	if err != nil {
		return err
	}

	if err = msg.Send(cli); err != nil {
		return err
	}

	return nil
}
