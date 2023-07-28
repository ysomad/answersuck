package round

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (s *Service) Update(ctx context.Context, nickname string, r entity.Round) error {
	p, err := s.pack.GetOne(ctx, r.PackID)
	if err != nil {
		return fmt.Errorf("s.pack.GetOne: %w", err)
	}

	if p.Author != nickname {
		return apperr.PackNotAuthor
	}

	return s.repo.UpdateOne(ctx, r)
}
