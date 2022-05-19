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

type ClientConfig struct {
	Host           string
	Port           int
	From           string
	Password       string
	KeepAlive      bool
	ConnectTimeout time.Duration
	SendTimeout    time.Duration
}

func (c *ClientConfig) Validate() error {
	switch {
	case c.Host == "":
		return ErrEmptyHost
	case c.From == "":
		return ErrEmptyFrom
	case c.Password == "":
		return ErrEmptyPassword
	case c.Port == 0:
		return ErrInvalidPort
	case c.ConnectTimeout == 0:
		c.ConnectTimeout = defaultConnectTimeout
	case c.SendTimeout == 0:
		c.SendTimeout = defaultSendTimeout
	}

	_, err := mail.ParseAddress(c.From)
	if err != nil {
		return err
	}

	return nil
}

func NewClient(cfg *ClientConfig) (*client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("c.Validate(): %w", err)
	}

	c := smail.NewSMTPClient()

	c.Host = cfg.Host
	c.Port = cfg.Port
	c.Username = cfg.From
	c.Password = cfg.Password
	c.Encryption = smail.EncryptionSSLTLS

	c.KeepAlive = cfg.KeepAlive
	c.ConnectTimeout = cfg.ConnectTimeout
	c.SendTimeout = cfg.SendTimeout

	return &client{
		srv:  c,
		from: cfg.From,
	}, nil
}

func (c *client) Send(ctx context.Context, e email.Email) error {
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
