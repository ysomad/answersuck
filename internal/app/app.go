package app

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/ysomad/answersuck/internal/config"
	"github.com/ysomad/answersuck/internal/pkg/httpserver"
)

func logFatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

type handlerContainer struct {
}

func Run(conf *config.Config, flags Flags) { //nolint:funlen // main func
	// pgClient, err := pgclient.New(
	// 	conf.PG.URL,
	// 	pgclient.WithMaxConns(conf.PG.MaxConns),
	// )
	// if err != nil {
	// 	logFatal("pgclient.New", err)
	// }

	// player
	// playerPG := player.NewPostgres(pgClient)
	// playerService := player.NewService(playerPG)

	// http
	mux := newServeMux("/rpc", handlerContainer{})

	srv := httpserver.New(mux, httpserver.WithPort(conf.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("received signal from httpserver", slog.String("signal", s.String()))
	case err := <-srv.Notify():
		slog.Info("got error from http server notify", slog.String("error", err.Error()))
	}

	if err := srv.Shutdown(); err != nil {
		slog.Info("got error on http server shutdown", slog.String("error", err.Error()))
	}
}
