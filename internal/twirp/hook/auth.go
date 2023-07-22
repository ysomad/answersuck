package hook

import (
	"context"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"golang.org/x/exp/slog"
)

type sessionService interface {
	Get(context.Context, string) (*session.Session, error)
}

func NewAuthorizedServerHooks(session sessionService) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			sid := appctx.GetSessionID(ctx)
			if sid == "" {
				return ctx, twirp.Unauthenticated.Error("session id not found in context")
			}

			s, err := session.Get(ctx, sid)
			if err != nil {
				slog.Error("couldnt get session", slog.String("error", err.Error()))
				return ctx, twirp.Unauthenticated.Error("session not found")
			}

			ctx = context.WithValue(ctx, appctx.SessionKey{}, s)

			return ctx, nil
		},
	}
}
