package round

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/service/common"
)

func (s *Service) Create(ctx context.Context, r entity.Round) (int32, error) {
	if err := common.VerifyAuthorship(ctx, s.pack, r.PackID); err != nil {
		return 0, fmt.Errorf("common.VerifyAuthorship: %w", err)
	}

	return s.repo.Save(ctx, r)
}
