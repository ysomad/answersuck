package hooks

import (
	"context"
	"log/slog"
	"time"

	"github.com/twitchtv/twirp"
)

type startTimeKey struct{}

func NewLogging() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, startTimeKey{}, time.Now()), nil
		},
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			s, _ := twirp.ServiceName(ctx)
			m, _ := twirp.MethodName(ctx)
			slog.Info("request routed", slog.String("service", s), slog.String("method", m))

			return ctx, nil
		},
		ResponseSent: func(ctx context.Context) {
			startTime := ctx.Value(startTimeKey{}).(time.Time)
			s, _ := twirp.ServiceName(ctx)
			m, _ := twirp.MethodName(ctx)
			slog.Info("response sent",
				slog.String("service", s),
				slog.String("method", m),
				slog.Duration("duration", time.Since(startTime)),
			)
		},
	}
}
