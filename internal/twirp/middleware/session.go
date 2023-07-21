package middleware

import (
	"context"
	"net/http"

	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"golang.org/x/exp/slog"
)

// WithSessionID writes session id from cookies into context.
func WithSessionID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sidCookie, err := r.Cookie("sid")
		if err != nil {
			slog.Debug("session id cookie not found")
			h.ServeHTTP(w, r)

			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, appctx.SessionIDKey, sidCookie.Value)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
