package account

import (
	"context"
	"fmt"
	"time"

	"github.com/ysomad/answersuck-backend/internal/config"
	"github.com/ysomad/answersuck-backend/internal/pkg/strings"
)

type service struct {
	cfg       *config.Aggregate
	repo      repository
	password  passwordManager
	session   sessionService
	email     emailService
	blockList blockList
}

type Deps struct {
	Config         *config.Aggregate
	AccountRepo    repository
	BlockList      blockList
	Password       passwordManager
	SessionService sessionService
	EmailService   emailService
}

func NewService(d *Deps) *service {
	return &service{
		cfg:       d.Config,
		repo:      d.AccountRepo,
		blockList: d.BlockList,
		password:  d.Password,
		session:   d.SessionService,
		email:     d.EmailService,
	}
}

func (s *service) Create(ctx context.Context, email, nickname, password string) (Account, error) {
	if s.blockList.Find(nickname) {
		return Account{}, fmt.Errorf("accountService - Create - s.blockList.Find: %w", ErrForbiddenNickname)
	}

	phash, err := s.password.Hash(password)
	if err != nil {
		return Account{}, fmt.Errorf("accountService - Create - s.password.Hash: %w", err)
	}

	code, err := strings.NewUnique(VerifCodeLen)
	if err != nil {
		return Account{}, fmt.Errorf("accountService - Create - a.generateVerifCode: %w", err)
	}

	now := time.Now()

	a, err := s.repo.Save(ctx, Account{
		Email:     email,
		Nickname:  nickname,
		Password:  phash,
		CreatedAt: now,
		UpdatedAt: now,
	}, code)
	if err != nil {
		return Account{}, fmt.Errorf("accountService - Create - s.repo.Save: %w", err)
	}

	go func() {
		_ = s.email.SendAccountVerificationEmail(ctx, a.Email, code)
	}()

	return a, nil
}

func (s *service) GetById(ctx context.Context, accountId string) (Account, error) {
	a, err := s.repo.FindById(ctx, accountId)
	if err != nil {
		return Account{}, fmt.Errorf("accountService - GetByID - s.repo.FindById: %w", err)
	}

	return a, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (Account, error) {
	a, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return Account{}, fmt.Errorf("accountService - GetByEmail - s.repo.FindByEmail: %w", err)
	}

	return a, nil
}

func (s *service) GetByNickname(ctx context.Context, nickname string) (Account, error) {
	a, err := s.repo.FindByNickname(ctx, nickname)
	if err != nil {
		return Account{}, fmt.Errorf("accountService - GetByNickname - s.repo.FindByNickname: %w", err)
	}

	return a, nil
}

func (s *service) Delete(ctx context.Context, accountId string) error {
	if err := s.repo.SetArchived(ctx, accountId, true, time.Now()); err != nil {
		return fmt.Errorf("accountService - Delete - s.repo.Archive: %w", err)
	}

	go func() {
		_ = s.session.TerminateAll(ctx, accountId)
	}()

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, accountId, oldPwd, newPwd string) error {
	oldHash, err := s.repo.FindPasswordById(ctx, accountId)
	if err != nil {
		return fmt.Errorf("accountService - UpdatePassword - s.repo.FindPasswordById: %w", err)
	}

	matches, err := s.password.Verify(oldPwd, oldHash)
	if err != nil {
		return fmt.Errorf("accountService - UpdatePassword - s.password.Verify: %w", err)
	}
	if !matches {
		return fmt.Errorf("accountService - UpdatePassword - s.password.Verify: %w", ErrInvalidPassword)
	}

	phash, err := s.password.Hash(newPwd)
	if err != nil {
		return fmt.Errorf("accountService - UpdatePassword - s.password.Hash: %w", err)
	}

	if err = s.repo.UpdatePassword(ctx, accountId, phash); err != nil {
		return fmt.Errorf("accountService - UpdatePassword - s.repo.UpdatePassword: %w", err)
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, login string) error {
	t, err := strings.NewUnique(PasswordTokenLen)
	if err != nil {
		return fmt.Errorf("accountService - ResetPassword - strings.NewUnique: %w", err)
	}

	email, err := s.repo.SavePasswordToken(ctx, SavePasswordTokenDTO{
		Login:     login,
		Token:     t,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("accountService - ResetPassword - s.repo.SavePasswordToken: %w", err)
	}

	go func() {
		_ = s.email.SendPasswordResetEmail(ctx, email, t)
	}()

	return nil
}

func (s *service) SetPassword(ctx context.Context, token, password string) error {
	t, err := s.repo.FindPasswordToken(ctx, token)
	if err != nil {
		return fmt.Errorf("accountService - SetPassword - s.repo.FindPasswordToken: %w", err)
	}

	if t.expired(s.cfg.Password.ResetTokenExpiration) {
		return fmt.Errorf("accountService - SetPassword - t.expired: %w", ErrPasswordTokenExpired)
	}

	phash, err := s.password.Hash(password)
	if err != nil {
		return fmt.Errorf("accountService - SetPassword - s.password.Hash: %w", err)
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
		_ = s.session.TerminateAll(ctx, t.AccountId)
	}()

	return nil
}

func (s *service) RequestVerification(ctx context.Context, accountId string) error {
	v, err := s.repo.FindVerification(ctx, accountId)
	if err != nil {
		return fmt.Errorf("accountService - RequestVerification - s.repo.FindVerification: %w", err)
	}

	if v.Verified {
		return fmt.Errorf("accountService - RequestVerification - v.Verified: %w", ErrAlreadyVerified)
	}

	go func() {
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
