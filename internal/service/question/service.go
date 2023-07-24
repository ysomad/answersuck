package question

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

type repository interface {
	Save(context.Context, *entity.Question) (int32, error)
}

type Service struct {
	repo repository
}

func NewService(r repository) *Service {
	return &Service{
		repo: r,
	}
}
