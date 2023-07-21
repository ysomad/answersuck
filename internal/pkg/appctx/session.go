package appctx

import (
	"context"

	"github.com/ysomad/answersuck/internal/pkg/session"
)

var (
	SessionIDKey = struct{}{}
	SessionKey   = struct{}{}
)

// GetSessionID returns sessions id from context or empty string if not found.
func GetSessionID(ctx context.Context) string {
	sid, ok := ctx.Value(SessionIDKey).(string)
	if !ok {
		return ""
	}

	return sid
}

// GetSession returns session from context or nil if session not found in
func GetSession(ctx context.Context) *session.Session {
	s, ok := ctx.Value(SessionKey).(*session.Session)
	if !ok {
		return nil
	}

	return s
}
