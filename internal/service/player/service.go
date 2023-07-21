package player

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

type repository interface {
	Save(context.Context, *entity.Player) error
	GetOne(ctx context.Context, login string, lt entity.LoginType) (*entity.Player, error)
}

type Service struct {
	repo repository
}

func NewService(r repository) *Service {
	return &Service{
		repo: r,
	}
}
