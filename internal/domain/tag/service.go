package tag

import (
	"context"
	"fmt"
)

type repository interface {
	SaveMultiple(ctx context.Context, r []CreateReq) ([]*Tag, error)
	FindAll(ctx context.Context) ([]*Tag, error)
}

type service struct {
	repo repository
}

func NewService(r repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateMultiple(ctx context.Context, r []CreateReq) ([]*Tag, error) {
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
