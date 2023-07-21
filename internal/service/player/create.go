package player

import (
	"context"
	"time"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/argon2"
)

func (s *Service) Create(ctx context.Context, nickname, email, pass string) error {
	hash, err := argon2.GenerateFromPassword(pass, argon2.DefaultParams)
	if err != nil {
		return err
	}

	return s.repo.Save(ctx, &entity.Player{
		Nickname:      nickname,
		Email:         email,
		EmailVerified: false,
		PasswordHash:  hash,
		CreatedAt:     time.Now(),
	})
}
