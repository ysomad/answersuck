package v1

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/session"
)

// mwAuthenticator check if request is authenticated and sets accountId and sessionId to locals (context)
func mwAuthenticator(l *zap.Logger, cfg *config.Session, s SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(cfg.CookieName)
			if err != nil {
				l.Info("http - v1 - middleware - mwAuthenticator - r.Cookie", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			sess, err := s.GetById(ctx, sessionCookie.Value)
			if err != nil {
				l.Error("http - v1 - middleware - mwAuthenticator - s.GetById", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if sess.Expired() {
				l.Info("http - v1 - middleware - mwAuthenticator", zap.Error(session.ErrExpired))
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
func mwVerificator(l *zap.Logger, cfg *config.Session, s SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(cfg.CookieName)
			if err != nil {
				l.Info("http - v1 - middleware - mwVerificator - r.Cookie", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			res, err := s.GetByIdWithDetails(ctx, sessionCookie.Value)
			if err != nil {
				l.Error("http - v1 - middleware - mwVerificator - s.GetById", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !res.Verified {
				l.Info("http - v1 - middleware - mwVerificator - !res.Verified", zap.Error(account.ErrNotEnoughRights))
				writeError(w, http.StatusForbidden, account.ErrNotEnoughRights)
				return
			}

			if res.Session.Expired() {
				l.Info("http - v1 - middleware - mwVerificator", zap.Error(session.ErrExpired))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !res.Session.SameDevice(r.RemoteAddr, r.UserAgent()) {
				l.Error("http - v1 - middleware - mwVerificator", zap.Error(session.ErrDeviceMismatch))
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
func mwTokenRequired(l *zap.Logger, token TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			accountId, err := getAccountId(r.Context())
			if err != nil {
				l.Info("http - v1 - middleware - mwTokenRequired - getAccountId", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			t := r.URL.Query().Get("token")
			if t == "" {
				l.Info("http - v1 - middleware - mwTokenRequired - r.URL.Query.Get", zap.Error(err))
				w.WriteHeader(http.StatusForbidden)
				return
			}

			sub, err := token.Parse(r.Context(), t)
			if err != nil {
				l.Error("http - v1 - middleware - mwTokenRequired - t.Parse", zap.Error(err))
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if sub != accountId {
				l.Info("http - v1 - middleware - mwTokenRequired - sub != accountId", zap.Error(err))
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
