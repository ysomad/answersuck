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
		GetMimeTypeById(ctx context.Context, mediaId string) (string, error)
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

func (s *service) Create(ctx context.Context, r CreateReq) (Answer, error) {
	a := Answer{Text: r.Text}

	if r.MediaId != "" {
		a.MediaId = &r.MediaId
	}

	if a.MediaId != nil {
		mimeType, err := s.media.GetMimeTypeById(ctx, r.MediaId)
		if err != nil {
			return Answer{}, fmt.Errorf("answerService - Create - s.media.GetMimeTypeById: %w", err)
		}

		if !a.isMimeTypeAllowed(mimeType) {
			return Answer{}, fmt.Errorf("answerService - Create - a.IsMimeTypeAllowed: %w", ErrMimeTypeNotAllowed)
		}
	}

	a, err := s.repo.Save(ctx, a)
	if err != nil {
		return Answer{}, fmt.Errorf("answerService - Create - s.repo.Save: %w", err)
	}

	return a, nil
}
