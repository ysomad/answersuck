package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type tagRepository interface {
	FindAll(ctx context.Context) ([]*domain.Tag, error)
}

type tag struct {
	repo tagRepository
}

func NewTag(r tagRepository) *tag {
	return &tag{
		repo: r,
	}
}

func (s *tag) GetAll(ctx context.Context) ([]*domain.Tag, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("tag - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
