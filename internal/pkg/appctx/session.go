package appctx

import (
	"context"

	"github.com/ysomad/answersuck/internal/pkg/session"
)

type (
	SessionIDKey struct{}
	SessionKey   struct{}
)

// GetSessionID returns sessions id from context or empty string if not found.
func GetSessionID(ctx context.Context) (string, bool) {
	sid, ok := ctx.Value(SessionIDKey{}).(string)
	return sid, ok
}

// GetSession returns session from context or nil if session not found in
func GetSession(ctx context.Context) (*session.Session, bool) {
	s, ok := ctx.Value(SessionKey{}).(*session.Session)
	return s, ok
}
