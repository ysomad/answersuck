package language

import (
	"context"
	"fmt"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Language, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) GetAll(ctx context.Context) ([]*Language, error) {
	l, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("languageService - GetAll - s.repo.FindAll: %w", err)
	}

	return l, nil
}
