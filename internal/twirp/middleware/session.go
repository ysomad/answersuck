package middleware

import (
	"context"
	"net"
	"net/http"

	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"golang.org/x/exp/slog"
)

// WithFootPrint writes remote addr and user agent into context.
func WithFootPrint(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(r.Header.Get("X-Real-IP"))
		if ip == nil {
			slog.Error("got invalid X-Real-IP header")
			h.ServeHTTP(w, r)

			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, appctx.FootPrintKey{}, appctx.FootPrint{
			IP:        ip,
			UserAgent: r.UserAgent(),
		})
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

// WithSessionID writes session id from cookies into context.
func WithSessionID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sidCookie, err := r.Cookie(session.Cookie)
		if err != nil {
			slog.Info("session id cookie not found")
			h.ServeHTTP(w, r)

			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, appctx.SessionIDKey{}, sidCookie.Value)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
