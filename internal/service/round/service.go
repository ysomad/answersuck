package round

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

type packService interface {
	GetOne(ctx context.Context, packID int32) (*entity.Pack, error)
}

type repository interface {
	Save(ctx context.Context, round entity.Round) (int32, error)
}

type Service struct {
	repo repository
	pack packService
}

func NewService(r repository, ps packService) *Service {
	return &Service{
		repo: r,
		pack: ps,
	}
}
