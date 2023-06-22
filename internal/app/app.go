package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/pgxpoolprometheus"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"

	"github.com/ysomad/answersuck/internal/config"
	"github.com/ysomad/answersuck/internal/pkg/httpserver"
	"github.com/ysomad/answersuck/internal/pkg/pgclient"
)

func logFatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func Run(conf *config.Config, flags Flags) { //nolint:funlen // main func
	otel, err := newOpenTelemetry(conf)
	if err != nil {
		logFatal("newOpenTelemetry", err)
	}

	defer otel.meterProvider.Shutdown(context.Background())
	defer otel.tracerProvider.Shutdown(context.Background())

	pgClient, err := pgclient.New(
		conf.PG.URL,
		pgclient.WithMaxConns(conf.PG.MaxConns),
		pgclient.WithQueryTracer(otel.pgxTracer),
	)
	if err != nil {
		logFatal("pgclient.New", err)
	}

	// pgx metrics
	pgxCollector := pgxpoolprometheus.NewCollector(pgClient.Pool, map[string]string{"db_name": conf.PG.DBName})
	if err = prometheus.Register(pgxCollector); err != nil {
		logFatal("prometheus.Register", err)
	}

	// http
	mux := newServeMux("/rpc")
	httpServer := httpserver.New(mux, httpserver.WithPort(conf.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("received signal from httpserver", slog.String("signal", s.String()))
	case err := <-httpServer.Notify():
		slog.Info("got error from http server notify", slog.String("error", err.Error()))
	}

	if err := httpServer.Shutdown(); err != nil {
		slog.Info("got error on http server shutdown", slog.String("error", err.Error()))
	}
}
