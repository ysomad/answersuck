package pack

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

// VerifyAuthorship returns no error if current user from session is author of pack.
func (s *Service) VerifyAuthorship(ctx context.Context, packID int32) error {
	nickname, ok := appctx.GetNickname(ctx)
	if !ok {
		return apperr.Unauthorized
	}

	pack, err := s.repo.GetOne(ctx, packID)
	if err != nil {
		return fmt.Errorf("error getting pack: %w", err)
	}

	if pack.Author != nickname {
		return apperr.PackNotAuthor
	}

	return nil
}

func (s *Service) VerifyRoundAuthorship(ctx context.Context, roundID int32) error {
	nickname, ok := appctx.GetNickname(ctx)
	if !ok {
		return apperr.Unauthorized
	}

	packAuthor, err := s.repo.GetRoundAuthor(ctx, roundID)
	if err != nil {
		return fmt.Errorf("error getting pack author: %w", err)
	}

	if packAuthor != nickname {
		return apperr.PackNotAuthor
	}

	return nil
}
