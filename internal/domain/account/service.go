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
	verifCodeLen = 64
	pwdTokenLen  = 64
)

type (
	Repository interface {
		Save(ctx context.Context, a *Account, code string) (string, error)
		FindEmailByNickname(ctx context.Context, nickname string) (string, error)
		FindById(ctx context.Context, accountId string) (*Account, error)
		FindByEmail(ctx context.Context, email string) (*Account, error)
		FindByNickname(ctx context.Context, nickname string) (*Account, error)
		Archive(ctx context.Context, accountId string, updatedAt time.Time) error

		Verify(ctx context.Context, code string, updatedAt time.Time) error
		FindVerification(ctx context.Context, nickname string) (VerificationDTO, error)

		SavePasswordToken(ctx context.Context, email, token string) error
		FindPasswordToken(ctx context.Context, token string) (PasswordToken, error)
		SetPassword(ctx context.Context, dto SetPasswordDTO) error
	}

	SessionService interface {
		TerminateAll(ctx context.Context, accountId string) error
	}

	EmailService interface {
		SendAccountVerificationEmail(ctx context.Context, to, code string) error
		SendPasswordResetEmail(ctx context.Context, to, token string) error
	}

	BlockList interface {
		Find(s string) bool
	}
)

type (
	service struct {
		cfg *config.Aggregate

		repo    Repository
		session SessionService
		email   EmailService

		blockList BlockList
	}

	Deps struct {
		Config         *config.Aggregate
		AccountRepo    Repository
		SessionService SessionService
		EmailService   EmailService
		BlockList      BlockList
	}
)

func NewService(d *Deps) *service {
	return &service{
		cfg:       d.Config,
		repo:      d.AccountRepo,
		session:   d.SessionService,
		email:     d.EmailService,
		blockList: d.BlockList,
	}
}

func (s *service) Create(ctx context.Context, r CreateRequest) (*Account, error) {
	now := time.Now()

	a := &Account{
		Email:     r.Email,
		Nickname:  r.Nickname,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if s.blockList.Find(a.Nickname) {
		return nil, fmt.Errorf("accountService - Create - s.blockList.Find: %w", ErrForbiddenNickname)
	}

	if err := a.setPassword(r.Password); err != nil {
		return nil, fmt.Errorf("accountService - Create - a.setPassword: %w", err)
	}

	c, err := a.generateVerifCode(verifCodeLen)
	if err != nil {
		return nil, fmt.Errorf("accountService - Create - a.generateVerifCode: %w", err)
	}

	accountId, err := s.repo.Save(ctx, a, c)
	if err != nil {
		return nil, fmt.Errorf("accountService - Create - s.repo.Save: %w", err)
	}

	a.Id = accountId

	go func() {
		// TODO: handle error
		_ = s.email.SendAccountVerificationEmail(ctx, a.Email, c)
	}()

	return a, nil
}

func (s *service) GetById(ctx context.Context, accountId string) (*Account, error) {
	a, err := s.repo.FindById(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("accountService - GetByID - s.repo.FindById: %w", err)
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

func (s *service) GetByNickname(ctx context.Context, nickname string) (*Account, error) {
	a, err := s.repo.FindByNickname(ctx, nickname)
	if err != nil {
		return nil, fmt.Errorf("accountService - GetByNickname - s.repo.FindByNickname: %w", err)
	}

	return a, nil
}

func (s *service) Delete(ctx context.Context, accountId string) error {
	if err := s.repo.Archive(ctx, accountId, time.Now()); err != nil {
		return fmt.Errorf("accountService - Delete - s.repo.Archive: %w", err)
	}

	go func() {
		// TODO: handle error
		_ = s.session.TerminateAll(ctx, accountId)
	}()

	return nil
}

func (s *service) RequestVerification(ctx context.Context, accountId string) error {
	v, err := s.repo.FindVerification(ctx, accountId)
	if err != nil {
		return fmt.Errorf("accountService - RequestVerification - s.repo.FindVerification: %w", err)
	}

	if v.Verified {
		return fmt.Errorf("accountService - v.Verified: %w", ErrAlreadyVerified)
	}

	go func() {
		// TODO: handle error
		_ = s.email.SendAccountVerificationEmail(ctx, v.Email, v.Code)
	}()

	return nil
}

func (s *service) Verify(ctx context.Context, code string) error {
	if err := s.repo.Verify(ctx, code, time.Now()); err != nil {
		return fmt.Errorf("accountService - Verify - s.repo.Verify: %w", err)
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, login string) error {
	email := login

	if _, err := mail.ParseAddress(login); err != nil {

		email, err = s.repo.FindEmailByNickname(ctx, login)
		if err != nil {
			return fmt.Errorf("accountService - ResetPassword - s.repo.FindEmailByNickname: %w", err)
		}

	}

	t, err := strings.NewUnique(pwdTokenLen)
	if err != nil {
		return fmt.Errorf("accountService - ResetPassword - strings.NewUnique: %w", err)
	}

	if err = s.repo.SavePasswordToken(ctx, email, t); err != nil {
		return fmt.Errorf("accountService - ResetPassword - s.repo.SavePasswordToken: %w", err)
	}

	go func() {
		// TODO: handle error
		_ = s.email.SendPasswordResetEmail(ctx, email, t)
	}()

	return nil
}

func (s *service) SetPassword(ctx context.Context, token, password string) error {
	t, err := s.repo.FindPasswordToken(ctx, token)
	if err != nil {
		return fmt.Errorf("accountService - SetPassword - s.repo.FindPasswordToken: %w", err)
	}

	if err = t.checkExpiration(s.cfg.Password.ResetTokenExp); err != nil {
		return fmt.Errorf("accountService - SetPassword - t.checkExpiration: %w", err)
	}

	phash, err := generatePasswordHash(password)
	if err != nil {
		return fmt.Errorf("accountService - SetPassword - a.generatePasswordHash: %w", err)
	}

	if err = s.repo.SetPassword(ctx, SetPasswordDTO{
		Token:     t.Token,
		AccountId: t.AccountId,
		Password:  phash,
		UpdatedAt: time.Now(),
	}); err != nil {
		return fmt.Errorf("accountService - SetPassword - s.repo.SetPassword: %w", err)
	}

	go func() {
		// TODO: handle error
		_ = s.session.TerminateAll(ctx, t.AccountId)
	}()

	return nil
}
