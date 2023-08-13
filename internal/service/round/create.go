package round

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

func (s *Service) Create(ctx context.Context, r entity.Round) (int32, error) {
	if err := s.pack.VerifyAuthorship(ctx, r.PackID); err != nil {
		return 0, err
	}

	return s.repo.Save(ctx, r)
}
