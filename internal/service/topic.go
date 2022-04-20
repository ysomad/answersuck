package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"
)

type topicService struct {
	repo TopicRepo
}

func NewTopicService(r TopicRepo) *topicService {
	return &topicService{
		repo: r,
	}
}

func (s *topicService) GetAll(ctx context.Context) ([]*domain.Topic, error) {
	t, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("topicService - GetAll - s.repo.FindAll: %w", err)
	}

	return t, nil
}
