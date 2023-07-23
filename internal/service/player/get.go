package player

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

func (s *Service) Get(ctx context.Context, login string) (*entity.Player, error) {
	if entity.NewLoginType(login) == entity.LoginTypeEmail {
		return s.repo.GetOne(ctx, login, entity.LoginTypeEmail)
	}

	return s.repo.GetOne(ctx, login, entity.LoginTypeNickname)
}
