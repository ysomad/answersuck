package answer

import (
	"context"
	"fmt"

	"github.com/answersuck/host/internal/pkg/pagination"
)

type (
	repository interface {
		Save(ctx context.Context, a Answer) (Answer, error)
		FindAll(ctx context.Context, p ListParams) (pagination.List[Answer], error)
	}

	mediaService interface {
		GetMediaTypeById(ctx context.Context, mediaId string) (string, error)
	}
)

type service struct {
	repo  repository
	media mediaService
}

func NewService(r repository, m mediaService) *service {
	return &service{
		repo:  r,
		media: m,
	}
}

func (s *service) Create(ctx context.Context, a Answer) (Answer, error) {
	if *a.MediaId != "" {
		mediaType, err := s.media.GetMediaTypeById(ctx, *a.MediaId)
		if err != nil {
			return Answer{}, fmt.Errorf("answerService - Create - s.media.GetMimeTypeById: %w", err)
		}

		if !mediaTypeAllowed(mediaType) {
			return Answer{}, fmt.Errorf("answerService - Create - mediaTypeAllowed: %w", ErrMediaTypeNotAllowed)
		}
	}

	a, err := s.repo.Save(ctx, a)
	if err != nil {
		return Answer{}, fmt.Errorf("answerService - Create - s.repo.Save: %w", err)
	}

	return a, nil
}

func (s *service) GetAll(ctx context.Context, p ListParams) (pagination.List[Answer], error) {
	answers, err := s.repo.FindAll(ctx, p)
	if err != nil {
		return pagination.List[Answer]{}, fmt.Errorf("answerService - GetAll - s.repo.FindAll: %w", err)
	}

	return answers, nil
}
