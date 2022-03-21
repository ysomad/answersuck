package service

import (
	"context"
	"fmt"

	"github.com/quizlyfun/quizly-backend/internal/config"

	"github.com/quizlyfun/quizly-backend/pkg/email"
)

const (
	verificationFormat = "%s/verification?code=%s"
)

type emailService struct {
	cfg    *config.Aggregate
	sender email.Sender
}

func NewEmailService(cfg *config.Aggregate, s email.Sender) *emailService {
	return &emailService{
		cfg:    cfg,
		sender: s,
	}
}

func (s *emailService) SendAccountVerificationEmail(ctx context.Context, to, username, code string) error {
	l := email.Letter{
		To:      to,
		Subject: fmt.Sprintf(s.cfg.Email.Subject.AccountVerification, username),
	}

	if err := l.SetBodyFromTemplate(
		s.cfg.Email.Template.AccountVerification,
		fmt.Sprintf(verificationFormat, s.cfg.Web.URL, code),
	); err != nil {
		return fmt.Errorf("emailService - SendEmailVerificationLetter - l.SetBodyFromTemplate: %w", err)
	}

	if err := s.sender.Send(l); err != nil {
		return fmt.Errorf("emailService - SendEmailVerificationLetter - s.sender.Send: %w", err)
	}

	return nil
}
