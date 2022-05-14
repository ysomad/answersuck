package service

import (
	"context"
	"fmt"
	"time"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"
)

type TopicRepository interface {
	Create(ctx context.Context, t domain.Topic) (int, error)
	FindAll(ctx context.Context) ([]*domain.Topic, error)
}

type topicService struct {
	repo TopicRepository
}

func NewTopicService(r TopicRepository) *topicService {
	return &topicService{
		repo: r,
	}
}

func (s *topicService) Create(ctx context.Context, req dto.TopicCreateRequest) (domain.Topic, error) {
	now := time.Now()

	t := domain.Topic{
		Name:       req.Name,
		LanguageId: req.LanguageId,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	topicId, err := s.repo.Create(ctx, t)
	if err != nil {
		return domain.Topic{}, fmt.Errorf("topicService - Create - s.repo.Create: %w", err)
	}

	t.Id = topicId

	return t, nil
}

func (s *topicService) GetAll(ctx context.Context) ([]*domain.Topic, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("topicService - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
