package service

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"

	"github.com/answersuck/vault/pkg/strings"
)

type accountPassword struct {
	cfg *config.Aggregate

	repo    AccountPasswordRepo
	session Session
	account Account
	email   Email
}

const (
	passwordResetTokenLength = 64
)

func NewAccountPasswordService(cfg *config.Aggregate, r AccountPasswordRepo, a Account, s Session, e Email) *accountPassword {
	return &accountPassword{
		cfg:     cfg,
		repo:    r,
		account: a,
		session: s,
		email:   e,
	}
}

func (s *accountPassword) RequestReset(ctx context.Context, login string) error {
	email := login

	if _, err := mail.ParseAddress(login); err != nil {

		a, err := s.account.GetByUsername(ctx, login)
		if err != nil {
			return fmt.Errorf("accountPasswordService - RequestReset - s.GetByUsername: %w", err)
		}

		email = a.Email
	}

	t, err := strings.NewUnique(passwordResetTokenLength)
	if err != nil {
		return fmt.Errorf("accountPasswordService - RequestReset - strings.NewUnique: %w", err)
	}

	if err = s.repo.InsertResetToken(ctx, email, t); err != nil {
		return fmt.Errorf("accountPasswordService - RequestReset - s.repo.InsertResetToken: %w", err)
	}

	if err = s.email.SendAccountPasswordResetMail(ctx, email, t); err != nil {
		return fmt.Errorf("accountPasswordService - RequestReset - s.email.SendAccountPasswordResetMail: %w", err)
	}

	return nil
}

func (s *accountPassword) Reset(ctx context.Context, token, password string) error {
	t, err := s.repo.FindResetToken(ctx, token)
	if err != nil {
		return fmt.Errorf("accountPasswordService - Reset - s.repo.FindResetToken: %w", err)
	}

	d := t.CreatedAt.Add(s.cfg.Password.ResetTokenExp)
	if time.Now().After(d) {
		return fmt.Errorf("accountPasswordService - Reset: %w", domain.ErrAccountResetPasswordTokenExpired)
	}

	a := domain.Account{Password: password}
	if err = a.GeneratePasswordHash(); err != nil {
		return fmt.Errorf("accountPasswordService - Reset - a.GeneratePassword: %w", err)
	}

	if err = s.repo.UpdateWithToken(ctx, dto.AccountUpdatePassword{
		Token:        t.Token,
		AccountId:    t.AccountId,
		PasswordHash: a.PasswordHash,
		UpdatedAt:    time.Now(),
	}); err != nil {
		return fmt.Errorf("accountPasswordService - Reset - s.repo.UpdatePasswordWithToken: %w", err)
	}

	return nil
}
