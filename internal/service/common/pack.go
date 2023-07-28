package common

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

type packService interface {
	GetOne(context.Context, int32) (*entity.Pack, error)
}

// VerifyAuthorship returns no error if current user from session is author of pack.
func VerifyAuthorship(ctx context.Context, s packService, packID int32) error {
	nickname, ok := appctx.GetNickname(ctx)
	if !ok {
		return apperr.Unauthorized
	}

	pack, err := s.GetOne(ctx, packID)
	if err != nil {
		return fmt.Errorf("s.pack.GetOne: %w", err)
	}

	if pack.Author != nickname {
		return apperr.PackNotAuthor
	}

	return nil
}
