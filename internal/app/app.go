package app

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/ysomad/answersuck/internal/config"
	mediapg "github.com/ysomad/answersuck/internal/pgrepo/media"
	playerpg "github.com/ysomad/answersuck/internal/pgrepo/player"
	questionpg "github.com/ysomad/answersuck/internal/pgrepo/question"
	tagpg "github.com/ysomad/answersuck/internal/pgrepo/tag"
	"github.com/ysomad/answersuck/internal/pkg/httpserver"
	"github.com/ysomad/answersuck/internal/pkg/pgclient"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"github.com/ysomad/answersuck/internal/service/auth"
	playersvc "github.com/ysomad/answersuck/internal/service/player"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	authtwirpv1 "github.com/ysomad/answersuck/internal/twirp/auth/v1"
	mediatwirpv1 "github.com/ysomad/answersuck/internal/twirp/media/v1"
	playertwirpv1 "github.com/ysomad/answersuck/internal/twirp/player/v1"
	questiontwirpv1 "github.com/ysomad/answersuck/internal/twirp/question/v1"
	tagtwirpv1 "github.com/ysomad/answersuck/internal/twirp/tag/v1"
)

func logFatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func Run(conf *config.Config, flags Flags) { //nolint:funlen // main func
	if flags.Migrate {
		mustRunMigrations(conf.PG.URL)
	}

	pgClient, err := pgclient.New(
		conf.PG.URL,
		pgclient.WithMaxConns(conf.PG.MaxConns),
	)
	if err != nil {
		logFatal("pgclient.New", err)
	}

	// session
	sessionPostgres := session.NewPostgresStore(pgClient.Pool)
	sessionManager := session.NewManager(sessionPostgres, conf.Session.LifeTime)

	// player
	playerPostgres := playerpg.NewRepository(pgClient)
	playerService := playersvc.NewService(playerPostgres)

	playerHandlerV1 := playertwirpv1.NewHandler(playerService)

	// tag
	tagPostgres := tagpg.NewRepository(pgClient)
	tagHandlerV1 := tagtwirpv1.NewHandler(tagPostgres, sessionPostgres)

	// auth
	authService := auth.NewService(sessionManager, playerService)
	authHandlerV1 := authtwirpv1.NewHandler(authService)

	// media
	mediaPostgres := mediapg.NewRepository(pgClient)
	mediaHandlerV1 := mediatwirpv1.NewHandler(mediaPostgres, sessionManager)

	// question
	questionPostgres := questionpg.NewRepository(pgClient)
	questionHandlerV1 := questiontwirpv1.NewHandler(questionPostgres, sessionManager)

	// http
	mux := apptwirp.NewMux([]apptwirp.Handler{
		playerHandlerV1,
		tagHandlerV1,
		authHandlerV1,
		mediaHandlerV1,
		questionHandlerV1,
	})

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
