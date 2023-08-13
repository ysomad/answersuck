package round

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

func (s *Service) Update(ctx context.Context, r entity.Round) error {
	if err := s.pack.VerifyAuthorship(ctx, r.PackID); err != nil {
		return err
	}

	return s.repo.UpdateOne(ctx, r)
}
