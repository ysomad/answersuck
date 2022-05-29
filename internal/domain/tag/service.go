package tag

import (
	"context"
	"fmt"
)

type Repository interface {
	SaveMultiple(ctx context.Context, r []CreateRequest) ([]*Tag, error)
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

func (s *service) CreateMultiple(ctx context.Context, r []CreateRequest) ([]*Tag, error) {
	t, err := s.repo.SaveMultiple(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("tagService - CreateMultiple - s.repo.SaveMultiple: %w", err)
	}

	return t, nil
}

func (s *service) GetAll(ctx context.Context) ([]*Tag, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("tagService - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
