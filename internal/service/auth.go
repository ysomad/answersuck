package service

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/quizlyfun/quizly-backend/internal/app"
	"github.com/quizlyfun/quizly-backend/internal/domain"

	"github.com/quizlyfun/quizly-backend/pkg/auth"
)

type authService struct {
	cfg     *app.Config
	token   auth.TokenManager
	account Account
	session Session
}

func NewAuthService(cfg *app.Config, t auth.TokenManager, a Account, s Session) *authService {
	return &authService{
		cfg:     cfg,
		token:   t,
		account: a,
		session: s,
	}
}

func (s *authService) Login(ctx context.Context, login, password string, d domain.Device) (domain.Session, error) {
	var a domain.Account

	_, err := mail.ParseAddress(login)
	if err != nil {
		// login is not email
		a, err = s.account.GetByUsername(ctx, login)
		if err != nil {
			return domain.Session{}, fmt.Errorf("authService - Login - s.account.GetByUsername: %w", err)
		}
	} else {
		// login is email
		a, err = s.account.GetByEmail(ctx, login)
		if err != nil {
			return domain.Session{}, fmt.Errorf("authService - Login - s.account.GetByEmail: %w", err)
		}
	}

	a.Password = password

	if err := a.CompareHashAndPassword(); err != nil {
		return domain.Session{}, fmt.Errorf("authService - Login - a.CompareHashAndPassword: %w", err)
	}

	sess, err := s.session.Create(ctx, a.Id, d)
	if err != nil {
		return domain.Session{}, fmt.Errorf("authService - Login - s.session.Create: %w", err)
	}

	return sess, nil
}

func (s *authService) Logout(ctx context.Context, sid string) error {
	if err := s.session.Terminate(ctx, sid); err != nil {
		return fmt.Errorf("authService - Logout - s.session.Terminate: %w", err)
	}

	return nil
}

func (s *authService) NewAccessToken(ctx context.Context, aid, password, audience string) (string, error) {
	a, err := s.account.GetByID(ctx, aid)
	if err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - s.account.GetByID: %w", err)
	}

	a.Password = password

	if err := a.CompareHashAndPassword(); err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - a.CompareHashAndPassword: %w", err)
	}

	t, err := s.token.New(aid, audience, s.cfg.AccessTokenTTL)
	if err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - s.token.New: %w", err)
	}

	return t, nil
}

func (s *authService) ParseAccessToken(ctx context.Context, token, audience string) (string, error) {
	aid, err := s.token.Parse(token, audience)
	if err != nil {
		return "", fmt.Errorf("authService - ParseAccessToken - s.token.Parse: %w", err)
	}

	return aid, nil
}
