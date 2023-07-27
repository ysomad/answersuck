package round

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (s *Service) Create(ctx context.Context, nickname string, r entity.Round) (int32, error) {
	p, err := s.pack.GetOne(ctx, r.PackID)
	if err != nil {
		return 0, fmt.Errorf("s.pack.GetOne: %w", err)
	}

	if p.Author != nickname {
		return 0, apperr.PackNotAuthor
	}

	return s.repo.Save(ctx, r)
}
