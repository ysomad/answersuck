package round

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/service/common"
)

func (s *Service) Update(ctx context.Context, r entity.Round) error {
	if err := common.VerifyPackAuthorship(ctx, s.pack, r.PackID); err != nil {
		return fmt.Errorf("common.VerifyPackAuthorship: %w", err)
	}

	return s.repo.UpdateOne(ctx, r)
}
