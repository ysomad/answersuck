package middleware

import (
	"context"
	"net/http"

	"github.com/ysomad/answersuck/internal/pkg/appctx"
)

// WithFootPrint writes remote addr and user agent into context.
func WithFootPrint(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, appctx.FootPrintKey, appctx.FootPrint{
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
		})
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
