package question

import (
	"context"
	"fmt"
	"time"
)

type Repository interface {
	Save(ctx context.Context, dto *CreateDTO) (int, error)
	FindById(ctx context.Context, questionId int) (*Question, error)
	FindAll(ctx context.Context) ([]Minimized, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) Create(ctx context.Context, dto *CreateDTO) (int, error) {
	now := time.Now()
	dto.CreatedAt = now
	dto.UpdatedAt = now

	questionId, err := s.repo.Save(ctx, dto)
	if err != nil {
		return 0, fmt.Errorf("questionService - Create - s.repo.Create: %w", err)
	}

	return questionId, nil
}

func (s *service) GetAll(ctx context.Context) ([]Minimized, error) {
	q, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("questionService - GetAll - s.repo.FindAll: %w", err)
	}

	return q, nil
}

func (s *service) GetById(ctx context.Context, questionId int) (*Question, error) {
	q, err := s.repo.FindById(ctx, questionId)
	if err != nil {
		return nil, fmt.Errorf("questionService - GetById - s.repo.FindById: %w", err)
	}

	return q, nil
}
