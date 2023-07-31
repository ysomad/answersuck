package middleware

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"golang.org/x/exp/slog"
)

type httpError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// WithFootPrint writes remote addr and user agent into context.
func WithFootPrint(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(r.Header.Get("X-Real-IP"))
		if ip == nil {
			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(httpError{
				Code: string(twirp.InvalidArgument),
				Msg:  apperr.MsgInvalidXRealIPHeader,
			}); err != nil {
				slog.Error("error encoding http error", slog.String("error", err.Error()))

				return
			}

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
