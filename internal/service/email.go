package service

import (
	"context"
	"fmt"

	"github.com/answersuck/answersuck-backend/internal/config"

	"github.com/answersuck/answersuck-backend/pkg/email"
)

const (
	verificationFormat = "%s/verification?code=%s"
)

type emailService struct {
	cfg   *config.Aggregate
	email email.Sender
}

type accountVerification struct {
	VerificationURL string
}

func NewEmailService(cfg *config.Aggregate, s email.Sender) *emailService {
	return &emailService{
		cfg:   cfg,
		email: s,
	}
}

func (s *emailService) SendAccountVerification(ctx context.Context, to, username, code string) error {
	l := email.Letter{
		To:      to,
		Subject: fmt.Sprintf(s.cfg.Email.Subject.AccountVerification, username),
	}

	if err := l.SetMsgFromTemplate(
		s.cfg.Email.Template.AccountVerification,
		accountVerification{fmt.Sprintf(verificationFormat, s.cfg.Web.URL, code)},
	); err != nil {
		return fmt.Errorf("emailService - SendAccountVerification - l.SetMsgFromTemplate: %w", err)
	}

	if err := s.email.Send(l); err != nil {
		return fmt.Errorf("emailService - SendAccountVerification - s.email.Send: %w", err)
	}

	return nil
}
