package question

import (
	"context"
	"fmt"
	"time"
)

type Repository interface {
	Save(ctx context.Context, dto *CreateDTO) (*Question, error)
	FindAll(ctx context.Context) ([]*Question, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) Create(ctx context.Context, dto *CreateDTO) (*Question, error) {
	now := time.Now()
	dto.CreatedAt = now
	dto.UpdatedAt = now

	q, err := s.repo.Save(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("questionService - Create - s.repo.Create: %w", err)
	}

	return q, nil
}

func (s *service) GetAll(ctx context.Context) ([]*Question, error) {
	q, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("questionService - GetAll - s.repo.FindAll: %w", err)
	}

	return q, nil
}
