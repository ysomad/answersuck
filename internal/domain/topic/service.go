package topic

import (
	"context"
	"fmt"
	"time"

	"github.com/ysomad/answersuck-backend/internal/pkg/pagination"
)

type repository interface {
	Save(ctx context.Context, t Topic) (Topic, error)
	FindAll(ctx context.Context, p ListParams) (pagination.List[Topic], error)
}

type service struct {
	repo repository
}

func NewService(r repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) Create(ctx context.Context, t Topic) (Topic, error) {
	t.CreatedAt = time.Now()

	t, err := s.repo.Save(ctx, t)
	if err != nil {
		return Topic{}, fmt.Errorf("topicService - Create - s.repo.Save: %w", err)
	}

	return t, nil
}

func (s *service) GetAll(ctx context.Context, p ListParams) (pagination.List[Topic], error) {
	tokenList, err := s.repo.FindAll(ctx, p)
	if err != nil {
		return pagination.List[Topic]{}, fmt.Errorf("topicService - GetAll - s.repo.FindAll: %w", err)
	}

	return tokenList, nil
}
