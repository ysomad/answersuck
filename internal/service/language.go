package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type LanguageRepository interface {
	FindAll(ctx context.Context) ([]*domain.Language, error)
}

type languageService struct {
	repo LanguageRepository
}

func NewLanguageService(r LanguageRepository) *languageService {
	return &languageService{
		repo: r,
	}
}

func (s *languageService) GetAll(ctx context.Context) ([]*domain.Language, error) {
	l, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("languageService - GetAll - s.repo.FindAll: %w", err)
	}

	return l, nil
}
