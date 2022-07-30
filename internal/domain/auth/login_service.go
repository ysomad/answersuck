package auth

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/session"
)

type loginService struct {
	account  AccountService
	session  SessionService
	password PasswordVerifier
}

func NewLoginService(a AccountService, s SessionService, p PasswordVerifier) *loginService {
	return &loginService{
		account:  a,
		session:  s,
		password: p,
	}
}

func (s *loginService) Login(ctx context.Context, login, password string, d session.Device) (*session.Session, error) {
	var a account.Account

	_, err := mail.ParseAddress(login)
	if err != nil {
		a, err = s.account.GetByNickname(ctx, login)
		if err != nil {
			return nil, fmt.Errorf("loginService - Login - s.account.GetByNickname: %w", err)
		}
	} else {
		a, err = s.account.GetByEmail(ctx, login)
		if err != nil {
			return nil, fmt.Errorf("loginService - Login - s.account.GetByEmail: %w", err)
		}
	}

	ok, err := s.password.Verify(password, a.Password)
	if err != nil {
		return nil, fmt.Errorf("tokenService - Create - s.password.Verify: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("tokenService - Create - s.password.Verify: %w", ErrIncorrectAccountPassword)
	}

	sess, err := s.session.Create(ctx, a.Id, d)
	if err != nil {
		return nil, fmt.Errorf("loginService - Login - s.session.Create: %w", err)
	}

	return sess, nil
}
