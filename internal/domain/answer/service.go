package answer

import (
	"context"
	"fmt"
)

type (
	repository interface {
		Save(ctx context.Context, a Answer) (Answer, error)
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

func (s *service) Create(ctx context.Context, text, mediaId string) (Answer, error) {
	if mediaId != "" {
		mediaType, err := s.media.GetMediaTypeById(ctx, mediaId)
		if err != nil {
			return Answer{}, fmt.Errorf("answerService - Create - s.media.GetMimeTypeById: %w", err)
		}

		if !mediaTypeAllowed(mediaType) {
			return Answer{}, fmt.Errorf("answerService - Create - a.IsMimeTypeAllowed: %w", ErrMediaTypeNotAllowed)
		}
	}

	a, err := s.repo.Save(ctx, Answer{
		Text:    text,
		MediaId: &mediaId,
	})
	if err != nil {
		return Answer{}, fmt.Errorf("answerService - Create - s.repo.Save: %w", err)
	}

	return a, nil
}
