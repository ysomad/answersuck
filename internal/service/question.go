package service

import (
	"context"
	"fmt"
	"time"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"
)

type QuestionRepository interface {
	Create(ctx context.Context, dto *dto.QuestionCreate) (*domain.Question, error)
	FindAll(ctx context.Context) ([]*domain.Question, error)
}

type questionService struct {
	repo QuestionRepository
}

func NewQuestionService(r QuestionRepository) *questionService {
	return &questionService{
		repo: r,
	}
}

func (s *questionService) Create(ctx context.Context, qc *dto.QuestionCreate) (*domain.Question, error) {
	now := time.Now()
	qc.CreatedAt = now
	qc.UpdatedAt = now

	q, err := s.repo.Create(ctx, qc)
	if err != nil {
		return nil, fmt.Errorf("questionService - Create - s.repo.Create: %w", err)
	}

	return q, nil
}

func (s *questionService) GetAll(ctx context.Context) ([]*domain.Question, error) {
	q, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("questionService - GetAll - s.repo.FindAll: %w", err)
	}

	return q, nil
}
