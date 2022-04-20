package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type tagService struct {
	repo TagRepo
}

func NewTagService(r TagRepo) *tagService {
	return &tagService{
		repo: r,
	}
}

func (s *tagService) GetAll(ctx context.Context) ([]*domain.Tag, error) {
	tags, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("tagService - GetAll - s.repo.FindAll: %w", err)
	}

	return tags, nil
}
