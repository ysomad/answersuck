package service

import (
	"context"
	"fmt"

	"github.com/quizlyfun/quizly-backend/internal/app"

	"github.com/quizlyfun/quizly-backend/pkg/email"
)

type emailService struct {
	cfg    *app.Config
	sender email.Sender
}

func NewEmailService(cfg *app.Config, s email.Sender) *emailService {
	return &emailService{
		cfg:    cfg,
		sender: s,
	}
}

func (s *emailService) SendEmailVerificationLetter(ctx context.Context, to, username, code string) error {
	l := email.Letter{
		To:      to,
		Subject: fmt.Sprintf(s.cfg.EmailVerificationSubject, username),
	}

	verifLink := fmt.Sprintf(s.cfg.EmailVerificationLink, code)

	if err := l.SetBodyFromTemplate(s.cfg.EmailVerificationTemplate, verifLink); err != nil {
		return fmt.Errorf("emailService - SendEmailVerificationLetter - l.SetBodyFromTemplate: %w", err)
	}

	if err := s.sender.Send(l); err != nil {
		return fmt.Errorf("emailService - SendEmailVerificationLetter - s.sender.Send: %w", err)
	}

	return nil
}
