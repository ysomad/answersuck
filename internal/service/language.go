package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type languageRepository interface {
	FindAll(ctx context.Context) ([]*domain.Language, error)
}

type language struct {
	repo languageRepository
}

func NewLanguage(r languageRepository) *language {
	return &language{
		repo: r,
	}
}

func (s *language) GetAll(ctx context.Context) ([]*domain.Language, error) {
	l, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("language - GetAll - s.repo.FindAll: %w", err)
	}

	return l, nil
}
