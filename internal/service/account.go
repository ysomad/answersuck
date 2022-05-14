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
	AccountRepository interface {
		Create(ctx context.Context, a *domain.Account) (*domain.Account, error)
		FindById(ctx context.Context, accountId string) (*domain.Account, error)
		FindByEmail(ctx context.Context, email string) (*domain.Account, error)
		FindByUsername(ctx context.Context, username string) (*domain.Account, error)
		Archive(ctx context.Context, accountId string, archived bool, updatedAt time.Time) error

		Verify(ctx context.Context, code string, verified bool, updatedAt time.Time) error
		FindVerification(ctx context.Context, accountId string) (dto.AccountVerification, error)

		InsertPasswordResetToken(ctx context.Context, email, token string) error
		FindPasswordResetToken(ctx context.Context, token string) (*dto.AccountPasswordResetToken, error)
		UpdatePasswordWithToken(ctx context.Context, a dto.AccountUpdatePassword) error
	}

	EmailService interface {
		SendAccountVerificationMail(ctx context.Context, to, code string) error
		SendAccountPasswordResetMail(ctx context.Context, to, token string) error
	}

	Finder interface {
		Find(s string) bool
	}
)

type accountService struct {
	cfg *config.Aggregate

	repo    AccountRepository
	session SessionService
	email   EmailService

	token     TokenManager
	blockList Finder
}

func NewAccountService(c *config.Aggregate, r AccountRepository, s SessionService,
	t TokenManager, e EmailService, f Finder) *accountService {
	return &accountService{
		cfg:       c,
		repo:      r,
		token:     t,
		session:   s,
		email:     e,
		blockList: f,
	}
}

func (s *accountService) Create(ctx context.Context, req dto.AccountCreateRequest) (*domain.Account, error) {
	now := time.Now()

	account := &domain.Account{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if s.blockList.Find(account.Username) {
		return nil, fmt.Errorf("accountService - Create - s.blockList.Find: %w", domain.ErrAccountForbiddenUsername)
	}

	if err := account.GeneratePasswordHash(); err != nil {
		return nil, fmt.Errorf("accountService - Create - acc.GeneratePasswordHash: %w", err)
	}

	account.SetDiceBearAvatar()

	if err := account.GenerateVerificationCode(verificationCodeLength); err != nil {
		return nil, fmt.Errorf("accountService - Create - a.GenerateVerificationCode: %w", err)
	}

	account, err := s.repo.Create(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("accountService- Create - s.repo.Create: %w", err)
	}

	go func() {
		_ = s.email.SendAccountVerificationMail(ctx, account.Email, account.VerificationCode)
	}()

	return account, nil
}

func (s *accountService) GetById(ctx context.Context, accountId string) (*domain.Account, error) {
	a, err := s.repo.FindById(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("accountService - GetByID - s.repo.FindByID: %w", err)
	}

	return a, nil
}

func (s *accountService) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	a, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("accountService - GetByEmail - s.repo.FindByEmail: %w", err)
	}

	return a, nil
}

func (s *accountService) GetByUsername(ctx context.Context, username string) (*domain.Account, error) {
	a, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("accountService - GetByUsername - s.repo.FindByUsername: %w", err)
	}

	return a, nil
}

func (s *accountService) Delete(ctx context.Context, accountId string) error {
	if err := s.repo.Archive(ctx, accountId, true, time.Now()); err != nil {
		return fmt.Errorf("accountService - Archive - s.repo.Archive: %w", err)
	}

	if err := s.session.TerminateAll(ctx, accountId); err != nil {
		return fmt.Errorf("accountService - Archive - s.session.TerminateAll: %w", err)
	}

	return nil
}

func (s *accountService) RequestVerification(ctx context.Context, accountId string) error {
	a, err := s.repo.FindVerification(ctx, accountId)
	if err != nil {
		return fmt.Errorf("accountService - RequestVerification - s.repo.FindById: %w", err)
	}

	if a.Verified {
		return fmt.Errorf("accountService: %w", domain.ErrAccountAlreadyVerified)
	}

	go func() {
		_ = s.email.SendAccountVerificationMail(ctx, a.Email, a.Code)
	}()

	return nil
}

func (s *accountService) Verify(ctx context.Context, code string, verified bool) error {
	if err := s.repo.Verify(ctx, code, verified, time.Now()); err != nil {
		return fmt.Errorf("accountService - Verify - s.repo.Verify: %w", err)
	}

	return nil
}

func (s *accountService) PasswordReset(ctx context.Context, login string) error {
	email := login

	if _, err := mail.ParseAddress(login); err != nil {

		a, err := s.GetByUsername(ctx, login)
		if err != nil {
			return fmt.Errorf("accountService - RequestPasswordReset - s.GetByUsername: %w", err)
		}

		email = a.Email
	}

	t, err := strings.NewUnique(passwordResetTokenLength)
	if err != nil {
		return fmt.Errorf("accountService - RequestPasswordReset - strings.NewUnique: %w", err)
	}

	if err = s.repo.InsertPasswordResetToken(ctx, email, t); err != nil {
		return fmt.Errorf("accountService - RequestPasswordReset - s.repo.InsertPasswordResetToken: %w", err)
	}

	if err = s.email.SendAccountPasswordResetMail(ctx, email, t); err != nil {
		return fmt.Errorf("accountService - RequestPasswordReset - s.email.SendAccountPasswordResetMail: %w", err)
	}

	return nil
}

func (s *accountService) PasswordSet(ctx context.Context, token, password string) error {
	t, err := s.repo.FindPasswordResetToken(ctx, token)
	if err != nil {
		return fmt.Errorf("accountService - PasswordReset - s.repo.FindPasswordResetToken: %w", err)
	}

	expiresAt := t.CreatedAt.Add(s.cfg.Password.ResetTokenExp)
	if time.Now().After(expiresAt) {
		return fmt.Errorf("accountService - PasswordReset: %w", domain.ErrAccountPasswordResetTokenExpired)
	}

	a := domain.Account{Password: password}

	if err = a.GeneratePasswordHash(); err != nil {
		return fmt.Errorf("accountService - PasswordReset - a.GeneratePassword: %w", err)
	}

	if err = s.repo.UpdatePasswordWithToken(ctx, dto.AccountUpdatePassword{
		Token:        t.Token,
		AccountId:    t.AccountId,
		PasswordHash: a.PasswordHash,
		UpdatedAt:    time.Now(),
	}); err != nil {
		return fmt.Errorf("accountService - PasswordReset - s.repo.UpdatePasswordWithToken: %w", err)
	}

	return nil
}
