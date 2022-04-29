package service

import (
	"context"
	"fmt"
	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"
	"time"
)

type topicRepository interface {
	Create(ctx context.Context, t domain.Topic) (int, error)
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

func (s *topic) Create(ctx context.Context, req dto.TopicCreateRequest) (domain.Topic, error) {
	now := time.Now()

	t := domain.Topic{
		Name:       req.Name,
		LanguageId: req.LanguageId,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	topicId, err := s.repo.Create(ctx, t)
	if err != nil {
		return domain.Topic{}, fmt.Errorf("topic - Create - s.repo.Create: %w", err)
	}

	t.Id = topicId

	return t, nil
}

func (s *topic) GetAll(ctx context.Context) ([]*domain.Topic, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("topic - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
