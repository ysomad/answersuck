package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/logging"
)

// mwAuthenticator check if request is authenticated and sets accountId and sessionId to locals (context)
func mwAuthenticator(l logging.Logger, cfg *config.Session, s SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(cfg.CookieName)
			if err != nil {
				l.Info("http - v1 - middleware - mwAuthenticator - r.Cookie :%w", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			sess, err := s.GetById(ctx, sessionCookie.Value)
			if err != nil {
				l.Error("http - v1 - middleware - mwAuthenticator - s.GetById: %w", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if sess.Expired() {
				l.Info("http - v1 - middleware - mwAuthenticator: %w", session.ErrExpired)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// if !sess.SameDevice(r.RemoteAddr, r.UserAgent()) {
			// 	l.Error("http - v1 - middleware - mwAuthenticator: %w", session.ErrDeviceMismatch)
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	return
			// }

			ctx = context.WithValue(ctx, sessionIdCtxKey{}, sess.Id)
			ctx = context.WithValue(ctx, accountIdCtxKey{}, sess.AccountId)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// mwVerificator is simillar to mwAuthenticator but also checks if account is verified,
// aborts if not.
//
// should be used instead of sessionMW
func mwVerificator(l logging.Logger, cfg *config.Session, s SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(cfg.CookieName)
			if err != nil {
				l.Info("http - v1 - middleware - mwVerificator - r.Cookie :%w", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			res, err := s.GetByIdWithVerified(ctx, sessionCookie.Value)
			if err != nil {
				l.Error("http - v1 - middleware - mwVerificator - s.GetById: %w", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !res.Verified {
				l.Info("http - v1 - middleware - mwVerificator - !res.Verified: %w", account.ErrNotEnoughRights)
				writeError(w, http.StatusForbidden, account.ErrNotEnoughRights)
				return
			}

			if time.Now().Unix() > res.Session.ExpiresAt {
				l.Info("http - v1 - middleware - mwVerificator: %w", session.ErrExpired)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if res.Session.IP != r.RemoteAddr || res.Session.UserAgent != r.UserAgent() {
				l.Error("http - v1 - middleware - mwVerificator: %w", session.ErrDeviceMismatch)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, sessionIdCtxKey{}, res.Session.Id)
			ctx = context.WithValue(ctx, accountIdCtxKey{}, res.Session.AccountId)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// mwTokenRequired parses and validates security token
func mwTokenRequired(l logging.Logger, token TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			accountId, err := getAccountId(r.Context())
			if err != nil {
				l.Info("http - v1 - middleware - mwTokenRequired - getAccountId: %w", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			t := r.URL.Query().Get("token")
			if t == "" {
				l.Info("http - v1 - middleware - mwTokenRequired - r.URL.Query.Get: %w", err)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			sub, err := token.Parse(r.Context(), t)
			if err != nil {
				l.Error("http - v1 - middleware - mwTokenRequired - t.Parse: %w", err)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if sub != accountId {
				l.Info("http - v1 - middleware - mwTokenRequired - sub!=accountId: %w", err)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// mwDeviceCtx sets session.Device object to context
func mwDeviceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), deviceCtxKey{}, session.Device{
			IP:        r.RemoteAddr,
			UserAgent: r.UserAgent(),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
