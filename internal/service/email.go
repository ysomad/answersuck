package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/pkg/email"
)

const (
	verificationFormat  = "%s/verify?code=%s"
	passwordResetFormat = "%s/reset?token=%s"
)

type emailService struct {
	cfg   *config.Aggregate
	email email.Sender
}

type withURL struct {
	URL string
}

func NewEmailService(cfg *config.Aggregate, s email.Sender) *emailService {
	return &emailService{
		cfg:   cfg,
		email: s,
	}
}

func (s *emailService) SendAccountVerificationMail(ctx context.Context, to, code string) error {
	if err := s.send(
		ctx,
		email.Letter{
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

func (s *emailService) SendAccountPasswordResetMail(ctx context.Context, to, token string) error {
	if err := s.send(
		ctx,
		email.Letter{
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

func (s *emailService) send(ctx context.Context, l email.Letter, tmplPath string, data interface{}) error {
	if err := l.SetMsgFromTemplate(tmplPath, data); err != nil {
		return fmt.Errorf("l.SetMsgFromTemplate: %w", err)
	}

	if err := s.email.Send(l); err != nil {
		return fmt.Errorf("s.email.Send: %w", err)
	}

	return nil
}
