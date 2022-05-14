package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type TagRepository interface {
	FindAll(ctx context.Context) ([]*domain.Tag, error)
}

type tagService struct {
	repo TagRepository
}

func NewTagService(r TagRepository) *tagService {
	return &tagService{
		repo: r,
	}
}

func (s *tagService) GetAll(ctx context.Context) ([]*domain.Tag, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("tagService - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
