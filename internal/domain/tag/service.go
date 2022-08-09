package tag

import (
	"context"
	"fmt"

	"github.com/answersuck/host/internal/pkg/pagination"
)

type repository interface {
	SaveMultiple(ctx context.Context, r []Tag) ([]Tag, error)
	FindAll(ctx context.Context, p ListParams) (pagination.List[Tag], error)
}

type service struct {
	repo repository
}

func NewService(r repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateMultiple(ctx context.Context, t []Tag) ([]Tag, error) {
	t, err := s.repo.SaveMultiple(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("tagService - CreateMultiple - s.repo.SaveMultiple: %w", err)
	}

	return t, nil
}

func (s *service) GetAll(ctx context.Context, p ListParams) (pagination.List[Tag], error) {
	l, err := s.repo.FindAll(ctx, p)
	if err != nil {
		return pagination.List[Tag]{}, fmt.Errorf("tagService - GetAll - s.repo.FindAll: %w", err)
	}

	return l, nil
}
