package v1

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/answersuck/host/internal/config"
	"github.com/answersuck/host/internal/domain/account"
	"github.com/answersuck/host/internal/domain/session"
)

// mwAuthenticator check if request is authenticated and sets accountId and sessionId to locals (context)
func mwAuthenticator(l *zap.Logger, cfg *config.Session, sess sessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(cfg.CookieName)
			if err != nil {
				l.Info("http - v1 - middleware - mwAuthenticator - r.Cookie", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			s, err := sess.GetById(ctx, sessionCookie.Value)
			if err != nil {
				l.Error("http - v1 - middleware - mwAuthenticator - s.GetById", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if s.Expired() {
				l.Info("http - v1 - middleware - mwAuthenticator", zap.Error(session.ErrExpired))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// if !sess.SameDevice(r.RemoteAddr, r.UserAgent()) {
			// 	l.Error("http - v1 - middleware - mwAuthenticator: %w", session.ErrDeviceMismatch)
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	return
			// }

			ctx = context.WithValue(ctx, sessionIdCtxKey{}, s.Id)
			ctx = context.WithValue(ctx, accountIdCtxKey{}, s.AccountId)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// mwVerificator is simillar to mwAuthenticator but also checks if account is verified,
// aborts if not.
//
// should be used instead of mwAuthenticator
func mwVerificator(l *zap.Logger, cfg *config.Session, sess sessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(cfg.CookieName)
			if err != nil {
				l.Info("http - v1 - middleware - mwVerificator - r.Cookie", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			s, err := sess.GetByIdWithDetails(ctx, sessionCookie.Value)
			if err != nil {
				l.Error("http - v1 - middleware - mwVerificator - s.GetById", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if s.Expired() {
				l.Info("http - v1 - middleware - mwVerificator", zap.Error(session.ErrExpired))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// if !sess.SameDevice(session.Device{IP: r.RemoteAddr, UserAgent: r.UserAgent()}) {
			// 	l.Error("http - v1 - middleware - mwVerificator", zap.Error(session.ErrDeviceMismatch))
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	return
			// }

			if !s.AccountVerified {
				l.Info("http - v1 - middleware - mwVerificator - !res.Verified", zap.Error(account.ErrNotEnoughRights))
				writeErr(w, http.StatusForbidden, account.ErrNotEnoughRights)
				return
			}

			ctx = context.WithValue(ctx, sessionIdCtxKey{}, s.Id)
			ctx = context.WithValue(ctx, accountIdCtxKey{}, s.AccountId)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// mwTokenRequired parses and validates security token
func mwTokenRequired(l *zap.Logger, token tokenService) func(http.Handler) http.Handler {
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
