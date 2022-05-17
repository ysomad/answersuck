package tag

import (
	"context"
	"fmt"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Tag, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) GetAll(ctx context.Context) ([]*Tag, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("tagService - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
