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

type service struct {
	cfg       *config.Aggregate
	repo      AccountRepo
	session   SessionService
	email     EmailService
	blockList BlockList
}

type Deps struct {
	Config         *config.Aggregate
	AccountRepo    AccountRepo
	SessionService SessionService
	EmailService   EmailService
	BlockList      BlockList
}

func NewService(d *Deps) *service {
	return &service{
		cfg:       d.Config,
		repo:      d.AccountRepo,
		session:   d.SessionService,
		email:     d.EmailService,
		blockList: d.BlockList,
	}
}

func (s *service) Create(ctx context.Context, r CreateReq) (*Account, error) {
	if s.blockList.Find(r.Nickname) {
		return nil, fmt.Errorf("accountService - Create - s.blockList.Find: %w", ErrForbiddenNickname)
	}

	now := time.Now()

	a := &Account{
		Email:     r.Email,
		Nickname:  r.Nickname,
		CreatedAt: now,
		UpdatedAt: now,
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

	var a Account

	phash, err := a.hashPassword(password)
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
