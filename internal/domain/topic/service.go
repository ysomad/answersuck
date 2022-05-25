package topic

import (
	"context"
	"fmt"
	"time"
)

type Repository interface {
	Save(ctx context.Context, t Topic) (Topic, error)
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

	t, err := s.repo.Save(ctx, t)
	if err != nil {
		return Topic{}, fmt.Errorf("topicService - Create - s.repo.Save: %w", err)
	}

	return t, nil
}

func (s *service) GetAll(ctx context.Context) ([]*Topic, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("topicService - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
