package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type questionService struct {
	repo QuestionRepo
}

func NewQuestionService(r QuestionRepo) *questionService {
	return &questionService{
		repo: r,
	}
}

func (s *questionService) GetAll(ctx context.Context) ([]*domain.Question, error) {
	q, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("questionService - GetAll - s.repo.FindAll: %w", err)
	}

	return q, nil
}
