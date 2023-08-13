package common

import (
	"context"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/session"
)

// CheckPlayerVerification verifies that request is authorized and player verified.
// Returns twirp errors.
func CheckPlayerVerification(ctx context.Context) (*session.Session, error) {
	s, ok := appctx.GetSession(ctx)
	if !ok {
		return nil, twirp.Unauthenticated.Error(apperr.MsgUnauthorized)
	}

	if !s.User.Verified {
		return nil, twirp.PermissionDenied.Error(apperr.MsgPlayerNotVerified)
	}

	return s, nil
}
