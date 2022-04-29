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

const (
	verificationCodeLength   = 64
	passwordResetTokenLength = 64
)

type (
	accountRepository interface {
		Create(ctx context.Context, a *domain.Account) (*domain.Account, error)
		FindById(ctx context.Context, aid string) (*domain.Account, error)
		FindByEmail(ctx context.Context, email string) (*domain.Account, error)
		FindByUsername(ctx context.Context, username string) (*domain.Account, error)
		Archive(ctx context.Context, aid string, archived bool, updatedAt time.Time) error

		Verify(ctx context.Context, code string, verified bool, updatedAt time.Time) error
		FindVerification(ctx context.Context, aid string) (dto.AccountVerification, error)

		InsertPasswordResetToken(ctx context.Context, email, token string) error
		FindPasswordResetToken(ctx context.Context, token string) (*dto.AccountPasswordResetToken, error)
		UpdatePasswordWithToken(ctx context.Context, dto dto.AccountUpdatePassword) error
	}

	emailService interface {
		SendAccountVerificationMail(ctx context.Context, to, code string) error
		SendAccountPasswordResetMail(ctx context.Context, to, token string) error
	}

	finder interface {
		Find(s string) bool
	}
)

type account struct {
	cfg *config.Aggregate

	repo    accountRepository
	session sessionService
	email   emailService

	token     tokenManager
	blockList finder
}

func NewAccount(c *config.Aggregate, r accountRepository, s sessionService,
	t tokenManager, e emailService, b finder) *account {
	return &account{
		cfg:       c,
		repo:      r,
		token:     t,
		session:   s,
		email:     e,
		blockList: b,
	}
}

func (s *account) Create(ctx context.Context, req dto.AccountCreateRequest) (*domain.Account, error) {
	a := &domain.Account{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if s.blockList.Find(a.Username) {
		return nil, fmt.Errorf("account - Create - s.blockList.Find: %w", domain.ErrAccountForbiddenUsername)
	}

	if err := a.GeneratePasswordHash(); err != nil {
		return nil, fmt.Errorf("account - Create - acc.GeneratePasswordHash: %w", err)
	}

	a.SetDiceBearAvatar()

	if err := a.GenerateVerificationCode(verificationCodeLength); err != nil {
		return nil, fmt.Errorf("account - Create - a.GenerateVerificationCode: %w", err)
	}

	a, err := s.repo.Create(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("account - Create - s.repo.Create: %w", err)
	}

	go func() {
		_ = s.email.SendAccountVerificationMail(ctx, a.Email, a.VerificationCode)
	}()

	return a, nil
}

func (s *account) GetById(ctx context.Context, aid string) (*domain.Account, error) {
	a, err := s.repo.FindById(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("account - GetByID - s.repo.FindByID: %w", err)
	}

	return a, nil
}

func (s *account) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	a, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("account - GetByEmail - s.repo.FindByEmail: %w", err)
	}

	return a, nil
}

func (s *account) GetByUsername(ctx context.Context, username string) (*domain.Account, error) {
	a, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("account - GetByUsername - s.repo.FindByUsername: %w", err)
	}

	return a, nil
}

func (s *account) Delete(ctx context.Context, aid, sid string) error {
	if err := s.repo.Archive(ctx, aid, true, time.Now()); err != nil {
		return fmt.Errorf("account - Archive - s.repo.Archive: %w", err)
	}

	if err := s.session.TerminateAll(ctx, aid); err != nil {
		return fmt.Errorf("account - Archive - s.session.TerminateAll: %w", err)
	}

	return nil
}

func (s *account) RequestVerification(ctx context.Context, aid string) error {
	a, err := s.repo.FindVerification(ctx, aid)
	if err != nil {
		return fmt.Errorf("account - RequestVerification - s.repo.FindById: %w", err)
	}

	if a.Verified {
		return fmt.Errorf("account: %w", domain.ErrAccountAlreadyVerified)
	}

	go func() {
		_ = s.email.SendAccountVerificationMail(ctx, a.Email, a.Code)
	}()

	return nil
}

func (s *account) Verify(ctx context.Context, code string, verified bool) error {
	if err := s.repo.Verify(ctx, code, verified, time.Now()); err != nil {
		return fmt.Errorf("account - Verify - s.repo.Verify: %w", err)
	}

	return nil
}

func (s *account) RequestPasswordReset(ctx context.Context, login string) error {
	accountEmail := login

	if _, err := mail.ParseAddress(login); err != nil {

		a, err := s.GetByUsername(ctx, login)
		if err != nil {
			return fmt.Errorf("account - RequestPasswordReset - s.GetByUsername: %w", err)
		}

		accountEmail = a.Email
	}

	t, err := strings.NewUnique(passwordResetTokenLength)
	if err != nil {
		return fmt.Errorf("account - RequestPasswordReset - strings.NewUnique: %w", err)
	}

	if err = s.repo.InsertPasswordResetToken(ctx, accountEmail, t); err != nil {
		return fmt.Errorf("account - RequestPasswordReset - s.repo.InsertPasswordResetToken: %w", err)
	}

	if err = s.email.SendAccountPasswordResetMail(ctx, accountEmail, t); err != nil {
		return fmt.Errorf("account - RequestPasswordReset - s.email.SendAccountPasswordResetMail: %w", err)
	}

	return nil
}

func (s *account) PasswordReset(ctx context.Context, token, password string) error {
	t, err := s.repo.FindPasswordResetToken(ctx, token)
	if err != nil {
		return fmt.Errorf("account - PasswordReset - s.repo.FindPasswordResetToken: %w", err)
	}

	d := t.CreatedAt.Add(s.cfg.Password.ResetTokenExp)
	if time.Now().After(d) {
		return fmt.Errorf("account - PasswordReset: %w", domain.ErrAccountResetPasswordTokenExpired)
	}

	a := domain.Account{Password: password}
	if err = a.GeneratePasswordHash(); err != nil {
		return fmt.Errorf("account - PasswordReset - a.GeneratePassword: %w", err)
	}

	if err = s.repo.UpdatePasswordWithToken(ctx, dto.AccountUpdatePassword{
		Token:        t.Token,
		AccountId:    t.AccountId,
		PasswordHash: a.PasswordHash,
		UpdatedAt:    time.Now(),
	}); err != nil {
		return fmt.Errorf("account - PasswordReset - s.repo.UpdatePasswordWithToken: %w", err)
	}

	return nil
}
