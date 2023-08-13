package hooks

import (
	"context"
	"log/slog"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/session"
)

func getSession(ctx context.Context, sm *session.Manager) (*session.Session, error) {
	sid, ok := appctx.GetSessionID(ctx)
	if !ok {
		return nil, twirp.Unauthenticated.Error("session id not found in context")
	}

	sess, err := sm.Get(ctx, sid)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// WithSession adds session to context if its present. Do not return error if
// player is not authenticated.
func WithSession(sm *session.Manager) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			sm, err := getSession(ctx, sm)
			if err != nil {
				slog.Info("error getting session", slog.String("error", err.Error()))
				return ctx, nil
			}

			ctx = context.WithValue(ctx, appctx.SessionKey{}, sm)
			ctx = context.WithValue(ctx, appctx.NicknameKey{}, sm.User.ID)

			return ctx, nil
		},
	}
}

// WithAuth authenticates player or returns unauthenticated error.
func WithAuth(sm *session.Manager) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			sess, err := getSession(ctx, sm)
			if err != nil {
				slog.Info("error getting session", slog.String("error", err.Error()))
				return ctx, twirp.Unauthenticated.Error(apperr.MsgUnauthorized)
			}

			ctx = context.WithValue(ctx, appctx.SessionKey{}, sess)
			ctx = context.WithValue(ctx, appctx.NicknameKey{}, sess.User.ID)

			return ctx, nil
		},
	}
}

// WithVerifiedPlayer authenticates player and checks whether he verified or not.
func WithVerifiedPlayer(sm *session.Manager) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			sess, err := getSession(ctx, sm)
			if err != nil {
				slog.Info("error getting session", slog.String("error", err.Error()))
				return ctx, twirp.Unauthenticated.Error(apperr.MsgUnauthorized)
			}

			if !sess.User.Verified {
				return ctx, twirp.PermissionDenied.Error(apperr.MsgPlayerNotVerified)
			}

			ctx = context.WithValue(ctx, appctx.SessionKey{}, sess)
			ctx = context.WithValue(ctx, appctx.NicknameKey{}, sess.User.ID)

			return ctx, nil
		},
	}
}
