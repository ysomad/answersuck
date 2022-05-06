package service

import (
	"context"
	"fmt"
	"time"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"
)

type questionRepository interface {
	Create(ctx context.Context, dto *dto.QuestionCreate) (*domain.Question, error)
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

func (s *question) Create(ctx context.Context, dto dto.QuestionCreate) (*domain.Question, error) {
	now := time.Now()
	dto.CreatedAt = now
	dto.UpdatedAt = now

	q, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("question - Create - s.repo.Create: %w", err)
	}

	return q, nil
}

func (s *question) GetAll(ctx context.Context) ([]*domain.Question, error) {
	q, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("question - GetAll - s.repo.FindAll: %w", err)
	}

	return q, nil
}
