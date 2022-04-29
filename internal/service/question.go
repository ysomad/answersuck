package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type questionRepository interface {
	FindAll(ctx context.Context) ([]*domain.Question, error)
}

type question struct {
	repo questionRepository
}

func NewQuestion(r questionRepository) *question {
	return &question{
		repo: r,
	}
}

func (s *question) GetAll(ctx context.Context) ([]*domain.Question, error) {
	q, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("question - GetAll - s.repo.FindAll: %w", err)
	}

	return q, nil
}
