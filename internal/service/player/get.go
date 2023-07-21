package player

import (
	"context"
	"net/mail"

	"github.com/ysomad/answersuck/internal/entity"
)

func (s *Service) Get(ctx context.Context, login string) (*entity.Player, error) {
	if _, err := mail.ParseAddress(login); err == nil {
		return s.repo.GetOne(ctx, login, entity.LoginTypeEmail)
	}

	return s.repo.GetOne(ctx, login, entity.LoginTypeNickname)
}
