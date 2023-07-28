package hooks

import (
	"context"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"golang.org/x/exp/slog"
)

type sessionGetter interface {
	Get(context.Context, string) (*session.Session, error)
}

func getSession(ctx context.Context, sg sessionGetter) (*session.Session, error) {
	sid, ok := appctx.GetSessionID(ctx)
	if !ok {
		return nil, twirp.Unauthenticated.Error("session id not found in context")
	}

	s, err := sg.Get(ctx, sid)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// NewSession adds session to context if its present.
func NewSession(sg sessionGetter) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			session, err := getSession(ctx, sg)
			if err != nil {
				slog.Info("error getting session", slog.String("error", err.Error()))
				return ctx, nil
			}

			ctx = context.WithValue(ctx, appctx.SessionKey{}, session)
			ctx = context.WithValue(ctx, appctx.NicknameKey{}, session.User.ID)

			return ctx, nil
		},
	}
}

// NewAuth returns Unauthenticated error if session id or session not found.
func NewAuth(sg sessionGetter) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			session, err := getSession(ctx, sg)
			if err != nil {
				slog.Info("error getting session", slog.String("error", err.Error()))
				return ctx, twirp.Unauthenticated.Error(apperr.MsgUnauthorized)
			}

			ctx = context.WithValue(ctx, appctx.SessionKey{}, session)
			ctx = context.WithValue(ctx, appctx.NicknameKey{}, session.User.ID)

			return ctx, nil
		},
	}
}

// NewAuthVerified returns Unauthenticated error if session id
// or session not found or PermissionDenied if player is not verified.
func NewAuthVerified(sg sessionGetter) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			session, err := getSession(ctx, sg)
			if err != nil {
				slog.Info("error getting session", slog.String("error", err.Error()))
				return ctx, twirp.Unauthenticated.Error(apperr.MsgUnauthorized)
			}

			if !session.User.Verified {
				return ctx, twirp.PermissionDenied.Error(apperr.MsgPlayerNotVerified)
			}

			ctx = context.WithValue(ctx, appctx.SessionKey{}, session)
			ctx = context.WithValue(ctx, appctx.NicknameKey{}, session.User.ID)

			return ctx, nil
		},
	}
}
