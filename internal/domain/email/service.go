package email

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/config"

	emailPkg "github.com/answersuck/vault/pkg/email"
)

//TODO: refactor

const (
	verificationFormat  = "%s/verify?code=%s"
	passwordResetFormat = "%s/reset?token=%s"
)

type Sender interface {
	Send(l emailPkg.Letter) error
}

type service struct {
	cfg   *config.Aggregate
	email Sender
}

type withURL struct {
	URL string
}

func NewService(cfg *config.Aggregate, s Sender) *service {
	return &service{
		cfg:   cfg,
		email: s,
	}
}

func (s *service) SendAccountVerificationMail(ctx context.Context, to, code string) error {
	if err := s.send(
		ctx,
		emailPkg.Letter{
			To:      to,
			Subject: s.cfg.Email.Subject.AccountVerification,
		},
		s.cfg.Email.Template.AccountVerification,
		withURL{fmt.Sprintf(verificationFormat, s.cfg.Web.URL, code)},
	); err != nil {
		return fmt.Errorf("emailService - SendAccountVerificationMail - s.send: %w", err)
	}

	return nil
}

func (s *service) SendAccountPasswordResetMail(ctx context.Context, to, token string) error {
	if err := s.send(
		ctx,
		emailPkg.Letter{
			To:      to,
			Subject: s.cfg.Email.Subject.AccountPasswordReset,
		},
		s.cfg.Email.Template.AccountPasswordReset,
		withURL{fmt.Sprintf(passwordResetFormat, s.cfg.Web.URL, token)},
	); err != nil {
		return fmt.Errorf("emailService - SendAccountVerificationMail - s.send: %w", err)
	}

	return nil
}

func (s *service) send(ctx context.Context, l emailPkg.Letter, tmplPath string, data interface{}) error {
	if err := l.SetMsgFromTemplate(tmplPath, data); err != nil {
		return fmt.Errorf("l.SetMsgFromTemplate: %w", err)
	}

	if err := s.email.Send(l); err != nil {
		return fmt.Errorf("s.email.Send: %w", err)
	}

	return nil
}
