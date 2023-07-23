package auth

import (
	"context"

	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/argon2"
	"github.com/ysomad/answersuck/internal/pkg/session"
)

func (s *Service) LogIn(ctx context.Context, login, password string, fp appctx.FootPrint) (*session.Session, error) {
	player, err := s.player.Get(ctx, login)
	if err != nil {
		return nil, err
	}

	ok, err := argon2.CompareHashAndPassword(password, player.PasswordHash)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, apperr.ErrInvalidCredentials
	}

	return s.session.Create(ctx, session.Player{
		Nickname:  player.Nickname,
		UserAgent: fp.UserAgent,
		IP:        fp.IP,
		Verified:  player.EmailVerified,
	})
}
