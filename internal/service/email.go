package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/config"

	emailPkg "github.com/answersuck/vault/pkg/email"
)

const (
	verificationFormat  = "%s/verify?code=%s"
	passwordResetFormat = "%s/reset?token=%s"
)

type emailSender interface {
	Send(l emailPkg.Letter) error
}

type email struct {
	cfg   *config.Aggregate
	email emailSender
}

type withURL struct {
	URL string
}

func NewEmail(cfg *config.Aggregate, s emailSender) *email {
	return &email{
		cfg:   cfg,
		email: s,
	}
}

func (s *email) SendAccountVerificationMail(ctx context.Context, to, code string) error {
	if err := s.send(
		ctx,
		emailPkg.Letter{
			To:      to,
			Subject: s.cfg.Email.Subject.AccountVerification,
		},
		s.cfg.Email.Template.AccountVerification,
		withURL{fmt.Sprintf(verificationFormat, s.cfg.Web.URL, code)},
	); err != nil {
		return fmt.Errorf("email - SendAccountVerificationMail - s.send: %w", err)
	}

	return nil
}

func (s *email) SendAccountPasswordResetMail(ctx context.Context, to, token string) error {
	if err := s.send(
		ctx,
		emailPkg.Letter{
			To:      to,
			Subject: s.cfg.Email.Subject.AccountPasswordReset,
		},
		s.cfg.Email.Template.AccountPasswordReset,
		withURL{fmt.Sprintf(passwordResetFormat, s.cfg.Web.URL, token)},
	); err != nil {
		return fmt.Errorf("email - SendAccountVerificationMail - s.send: %w", err)
	}

	return nil
}

func (s *email) send(ctx context.Context, l emailPkg.Letter, tmplPath string, data interface{}) error {
	if err := l.SetMsgFromTemplate(tmplPath, data); err != nil {
		return fmt.Errorf("l.SetMsgFromTemplate: %w", err)
	}

	if err := s.email.Send(l); err != nil {
		return fmt.Errorf("s.email.Send: %w", err)
	}

	return nil
}
