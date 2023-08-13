package roundquestion

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

type repository interface {
	Save(ctx context.Context, round *entity.RoundQuestion) (int32, error)
}

type Service struct {
	repo repository
}

func NewService(r repository) *Service {
	return &Service{
		repo: r,
	}
}
