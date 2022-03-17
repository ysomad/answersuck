package service

import (
	"context"
	"fmt"

	"github.com/quizlyfun/quizly-backend/internal/app"
	"github.com/quizlyfun/quizly-backend/pkg/email"
)

type emailService struct {
	cfg *app.Config
}

func NewEmailService(cfg *app.Config) *emailService {
	return &emailService{
		cfg: cfg,
	}
}

func (s *emailService) SendAccountVerificationEmail(ctx context.Context, to string) error {
	_ = email.Letter{
		To:      to,
		Subject: "",
	}

	_, err := email.NewClient(s.cfg.SMTPFrom, s.cfg.SMTPPass, s.cfg.SMTPHost, s.cfg.SMTPPort)
	if err != nil {
		return fmt.Errorf("emailService - SendAccountVerification - s.repo.Create: %w", err)
	}

	panic("implement")

	return nil
}
