package pack

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

type repository interface {
	GetOne(context.Context, int32) (*entity.Pack, error)
	GetRoundAuthor(ctx context.Context, roundID int32) (string, error)
}

type Service struct {
	repo repository
}

func NewService(r repository) *Service {
	return &Service{
		repo: r,
	}
}
