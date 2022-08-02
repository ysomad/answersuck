package email

import (
	"context"
	"fmt"

	"github.com/answersuck/host/internal/config"
)

type smtpClient interface {
	SendEmail(ctx context.Context, e Email) error
}

type service struct {
	cfg    *config.Email
	webURL string
	client smtpClient
}

func NewService(cfg *config.Aggregate, s smtpClient) *service {
	return &service{
		cfg:    &cfg.Email,
		webURL: cfg.Web.URL,
		client: s,
	}
}

func (s *service) SendAccountVerificationEmail(ctx context.Context, to, code string) error {
	return s.send(ctx, sendEmailDTO{
		to:         to,
		template:   s.cfg.Template.AccountVerification,
		subject:    s.cfg.Subject.AccountVerification,
		format:     s.cfg.Format.AccountVerification,
		formatArgs: []any{s.webURL, code},
	})
}

func (s *service) SendPasswordResetEmail(ctx context.Context, to, token string) error {
	return s.send(ctx, sendEmailDTO{
		to:         to,
		template:   s.cfg.Template.AccountPasswordReset,
		subject:    s.cfg.Subject.AccountPasswordReset,
		format:     s.cfg.Format.AccountPasswordReset,
		formatArgs: []any{s.webURL, token},
	})
}

// Private

func (s *service) send(ctx context.Context, dto sendEmailDTO) error {
	e := Email{
		To:      dto.to,
		Subject: dto.subject,
	}

	url := fmt.Sprintf(dto.format, dto.formatArgs...)
	err := e.setMessageFromTemplate(dto.template, withURL{url})
	if err != nil {
		return fmt.Errorf("e.SetMessageFromTemplate: %w", err)
	}

	if err := e.validate(); err != nil {
		return fmt.Errorf("e.validate: %w", err)
	}

	if err := s.client.SendEmail(ctx, e); err != nil {
		return fmt.Errorf("s.client.SendEmail: %w", err)
	}

	return nil
}
