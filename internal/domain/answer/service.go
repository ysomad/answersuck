package answer

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/pkg/logging"
)

type (
	Repository interface {
		Save(ctx context.Context, a Answer) (Answer, error)
	}

	MediaService interface {
		GetMimeTypeById(ctx context.Context, mediaId string) (string, error)
	}
)

type service struct {
	log   logging.Logger
	repo  Repository
	media MediaService
}

func NewService(l logging.Logger, r Repository, m MediaService) *service {
	return &service{
		log:   l,
		repo:  r,
		media: m,
	}
}

func (s *service) Create(ctx context.Context, r CreateRequest) (Answer, error) {
	mimeType, err := s.media.GetMimeTypeById(ctx, r.MediaId)
	if err != nil {
		return Answer{}, fmt.Errorf("answerService - Create - s.media.FindMimeTypeById: %w", err)
	}

	s.log.Info(mimeType)

	a := Answer{
		Text:    r.Text,
		MediaId: r.MediaId,
	}

	if !a.isMimeTypeAllowed(mimeType) {
		return Answer{}, fmt.Errorf("answerService - Create - a.IsMimeTypeAllowed: %w", ErrMimeTypeNotAllowed)
	}

	a, err = s.repo.Save(ctx, a)
	if err != nil {
		return Answer{}, fmt.Errorf("answerService - Create - s.repo.Create: %w", err)
	}

	return a, nil
}
