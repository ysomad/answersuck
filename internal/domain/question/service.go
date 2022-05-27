package question

import (
	"context"
	"fmt"
)

type Repository interface {
	Save(ctx context.Context, q *Question) (int, error)
	FindById(ctx context.Context, questionId int) (*Detailed, error)
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

func (s *service) Create(ctx context.Context, q *Question) (*Question, error) {
	q.PrepareForSave()

	questionId, err := s.repo.Save(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("questionService - Create - s.repo.Save: %w", err)
	}

	q.Id = questionId

	return q, nil
}

func (s *service) GetAll(ctx context.Context) ([]Minimized, error) {
	q, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("questionService - GetAll - s.repo.FindAll: %w", err)
	}

	return q, nil
}

func (s *service) GetById(ctx context.Context, questionId int) (*Detailed, error) {
	q, err := s.repo.FindById(ctx, questionId)
	if err != nil {
		return nil, fmt.Errorf("questionService - GetById - s.repo.FindById: %w", err)
	}

	return q, nil
}
