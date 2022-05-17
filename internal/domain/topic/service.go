package topic

import (
	"context"
	"fmt"
	"time"
)

type Repository interface {
	Create(ctx context.Context, t Topic) (int, error)
	FindAll(ctx context.Context) ([]*Topic, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) Create(ctx context.Context, r CreateRequest) (Topic, error) {
	now := time.Now()

	t := Topic{
		Name:       r.Name,
		LanguageId: r.LanguageId,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	topicId, err := s.repo.Create(ctx, t)
	if err != nil {
		return Topic{}, fmt.Errorf("topicService - Create - s.repo.Create: %w", err)
	}

	t.Id = topicId

	return t, nil
}

func (s *service) GetAll(ctx context.Context) ([]*Topic, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("topicService - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
