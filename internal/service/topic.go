package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type topicRepository interface {
	FindAll(ctx context.Context) ([]*domain.Topic, error)
}

type topic struct {
	repo topicRepository
}

func NewTopic(r topicRepository) *topic {
	return &topic{
		repo: r,
	}
}

func (s *topic) GetAll(ctx context.Context) ([]*domain.Topic, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("topic - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
