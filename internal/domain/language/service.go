package language

import (
	"context"
	"fmt"
)

type repository interface {
	FindAll(ctx context.Context) ([]Language, error)
}

type service struct {
	repo repository
}

func NewService(r repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) GetAll(ctx context.Context) ([]Language, error) {
	l, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("languageService - GetAll - s.repo.FindAll: %w", err)
	}

	return l, nil
}
