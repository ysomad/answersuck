package account

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/answersuck/vault/internal/config"

	"github.com/answersuck/vault/pkg/strings"
)

const (
	verificationCodeLength   = 64
	passwordResetTokenLength = 64
)

type (
	Repository interface {
		Save(ctx context.Context, a *Account, code string) (string, error)
		FindById(ctx context.Context, accountId string) (*Account, error)
		FindByEmail(ctx context.Context, email string) (*Account, error)
		FindByNickname(ctx context.Context, nickname string) (*Account, error)
		SetArchived(ctx context.Context, accountId string, archived bool, updatedAt time.Time) error

		Verify(ctx context.Context, code string, verified bool, updatedAt time.Time) error
		FindVerification(ctx context.Context, accountId string) (VerificationDTO, error)

		SavePasswordResetToken(ctx context.Context, email, token string) error
		FindPasswordResetToken(ctx context.Context, token string) (PasswordResetToken, error)
		SetPassword(ctx context.Context, dto SetPasswordDTO) error
	}

	SessionService interface {
		TerminateAll(ctx context.Context, accountId string) error
	}

	EmailService interface {
		SendAccountVerificationEmail(ctx context.Context, to, code string) error
		SendPasswordResetEmail(ctx context.Context, to, token string) error
	}

	Blocklist interface {
		Find(s string) bool
	}
)

type service struct {
	cfg *config.Aggregate

	repo    Repository
	session SessionService
	email   EmailService

	blocklist Blocklist
}

type Deps struct {
	Config         *config.Aggregate
	AccountRepo    Repository
	SessionService SessionService
	EmailService   EmailService
	Blocklist      Blocklist
}

func NewService(d *Deps) *service {
	return &service{
		cfg:       d.Config,
		repo:      d.AccountRepo,
		session:   d.SessionService,
		email:     d.EmailService,
		blocklist: d.Blocklist,
	}
}

func (s *service) Create(ctx context.Context, r CreateRequest) (*Account, error) {
	now := time.Now()

	a := &Account{
		Email:     r.Email,
		Nickname:  r.Nickname,
		Password:  r.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if s.blocklist.Find(a.Nickname) {
		return nil, fmt.Errorf("accountService - Create - s.blockList.Find: %w", ErrForbiddenNickname)
	}

	if err := a.generatePasswordHash(); err != nil {
		return nil, fmt.Errorf("accountService - Create - a.generatePasswordHash: %w", err)
	}

	c, err := a.generateVerificationCode(verificationCodeLength)
	if err != nil {
		return nil, fmt.Errorf("accountService - Create - a.generateVerificationCode: %w", err)
	}

	a, err = s.repo.Save(ctx, a, c)
	if err != nil {
		return nil, fmt.Errorf("accountService - Create - s.repo.Save: %w", err)
	}

	go func() {
		// TODO: handle error
		_ = s.email.SendAccountVerificationEmail(ctx, a.Email, c)
	}()

	return a, nil
}

func (s *service) GetById(ctx context.Context, accountId string) (*Account, error) {
	a, err := s.repo.FindById(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("accountService - GetByID - s.repo.FindByID: %w", err)
	}

	return a, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*Account, error) {
	a, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("accountService - GetByEmail - s.repo.FindByEmail: %w", err)
	}

	return a, nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (*Account, error) {
	a, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("accountService - GetByUsername - s.repo.FindByUsername: %w", err)
	}

	return a, nil
}

func (s *service) Delete(ctx context.Context, accountId string) error {
	if err := s.repo.SetArchived(ctx, accountId, true, time.Now()); err != nil {
		return fmt.Errorf("accountService - Archive - s.repo.SetArchived: %w", err)
	}

	if err := s.session.TerminateAll(ctx, accountId); err != nil {
		return fmt.Errorf("accountService - Archive - s.session.TerminateAll: %w", err)
	}

	return nil
}

func (s *service) RequestVerification(ctx context.Context, accountId string) error {
	v, err := s.repo.FindVerification(ctx, accountId)
	if err != nil {
		return fmt.Errorf("accountService - RequestVerification - s.repo.FindVerification: %w", err)
	}

	if v.Verified {
		return fmt.Errorf("accountService: %w", ErrAlreadyVerified)
	}

	go func() {
		// TODO: handle error
		_ = s.email.SendAccountVerificationEmail(ctx, v.Email, v.Code)
	}()

	return nil
}

func (s *service) Verify(ctx context.Context, code string, verified bool) error {
	if err := s.repo.Verify(ctx, code, verified, time.Now()); err != nil {
		return fmt.Errorf("accountService - Verify - s.repo.Verify: %w", err)
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, login string) error {
	email := login

	if _, err := mail.ParseAddress(login); err != nil {

		a, err := s.GetByUsername(ctx, login)
		if err != nil {
			return fmt.Errorf("accountService - ResetPassword - s.GetByUsername: %w", err)
		}

		email = a.Email
	}

	t, err := strings.NewUnique(passwordResetTokenLength)
	if err != nil {
		return fmt.Errorf("accountService - ResetPassword - strings.NewUnique: %w", err)
	}

	if err = s.repo.SavePasswordResetToken(ctx, email, t); err != nil {
		return fmt.Errorf("accountService - ResetPassword - s.repo.SavePasswordResetToken: %w", err)
	}

	go func() {
		// TODO: handle error
		_ = s.email.SendPasswordResetEmail(ctx, email, t)
	}()

	return nil
}

func (s *service) SetPassword(ctx context.Context, token, password string) error {
	t, err := s.repo.FindPasswordResetToken(ctx, token)
	if err != nil {
		return fmt.Errorf("accountService - SetPassword - s.repo.FindPasswordResetToken: %w", err)
	}

	if err = t.checkExpiration(s.cfg.Password.ResetTokenExp); err != nil {
		return fmt.Errorf("accountService - SetPassword - t.CheckExpired: %w", err)
	}

	a := Account{Password: password}

	if err = a.generatePasswordHash(); err != nil {
		return fmt.Errorf("accountService - SetPassword - a.GeneratePassword: %w", err)
	}

	if err = s.repo.SetPassword(ctx, SetPasswordDTO{
		Token:        t.Token,
		AccountId:    t.AccountId,
		PasswordHash: a.PasswordHash,
		UpdatedAt:    time.Now(),
	}); err != nil {
		return fmt.Errorf("accountService - SetPassword - s.repo.SetPassword: %w", err)
	}

	go func() {
		// TODO: handle error
		_ = s.session.TerminateAll(ctx, t.AccountId)
	}()

	return nil
}
