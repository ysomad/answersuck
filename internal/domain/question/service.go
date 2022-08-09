package question

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type repository interface {
	Save(ctx context.Context, dto CreateDTO) (questionId uint32, err error)
	FindById(ctx context.Context, questionId uint32) (Detailed, error)
	FindAll(ctx context.Context) ([]Minimized, error)
}

type mediaProvider interface {
	URL(filename string) *url.URL
}

type service struct {
	repo  repository
	media mediaProvider
}

func NewService(r repository, p mediaProvider) *service {
	return &service{
		repo:  r,
		media: p,
	}
}

func (s *service) Create(ctx context.Context, dto CreateDTO) (uint32, error) {
	dto.CreatedAt = time.Now()

	questionId, err := s.repo.Save(ctx, dto)
	if err != nil {
		return 0, fmt.Errorf("questionService - Create - s.repo.Save: %w", err)
	}

	return questionId, nil
}

func (s *service) GetById(ctx context.Context, questionId uint32) (Detailed, error) {
	d, err := s.repo.FindById(ctx, questionId)
	if err != nil {
		return Detailed{}, fmt.Errorf("questionService - GetById - s.repo.FindById: %w", err)
	}

	d.setURLsFromFilenames(s.media)

	return d, nil
}

func (s *service) GetAll(ctx context.Context) ([]Minimized, error) {
	q, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("questionService - GetAll - s.repo.FindAll: %w", err)
	}

	return q, nil
}
