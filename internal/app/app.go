package app

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/ysomad/answersuck/internal/config"

	mediapg "github.com/ysomad/answersuck/internal/pgrepo/media"
	packpg "github.com/ysomad/answersuck/internal/pgrepo/pack"
	playerpg "github.com/ysomad/answersuck/internal/pgrepo/player"
	questionpg "github.com/ysomad/answersuck/internal/pgrepo/question"
	roundpg "github.com/ysomad/answersuck/internal/pgrepo/round"
	tagpg "github.com/ysomad/answersuck/internal/pgrepo/tag"

	authsvc "github.com/ysomad/answersuck/internal/service/auth"
	playersvc "github.com/ysomad/answersuck/internal/service/player"
	roundsvc "github.com/ysomad/answersuck/internal/service/round"

	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	authv1 "github.com/ysomad/answersuck/internal/twirp/auth/v1"
	mediav1 "github.com/ysomad/answersuck/internal/twirp/editor/v1/media"
	packv1 "github.com/ysomad/answersuck/internal/twirp/editor/v1/pack"
	questionv1 "github.com/ysomad/answersuck/internal/twirp/editor/v1/question"
	roundv1 "github.com/ysomad/answersuck/internal/twirp/editor/v1/round"
	tagv1 "github.com/ysomad/answersuck/internal/twirp/editor/v1/tag"
	playerv1 "github.com/ysomad/answersuck/internal/twirp/player/v1"

	"github.com/ysomad/answersuck/internal/pkg/httpserver"
	"github.com/ysomad/answersuck/internal/pkg/pgclient"
	"github.com/ysomad/answersuck/internal/pkg/session"
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

	playerHandlerV1 := playerv1.NewHandler(playerService)

	// tag
	tagPostgres := tagpg.NewRepository(pgClient)
	tagHandlerV1 := tagv1.NewHandler(tagPostgres, sessionPostgres)

	// auth
	authService := authsvc.NewService(sessionManager, playerService)
	authHandlerV1 := authv1.NewHandler(authService)

	// media
	mediaPostgres := mediapg.NewRepository(pgClient)
	mediaHandlerV1 := mediav1.NewHandler(mediaPostgres, sessionManager)

	// question
	questionPostgres := questionpg.NewRepository(pgClient)
	questionHandlerV1 := questionv1.NewHandler(questionPostgres, sessionManager)

	// pack
	packPostgres := packpg.NewRepository(pgClient)
	packHandlerV1 := packv1.NewHandler(packPostgres, sessionManager)

	// round
	roundPostgres := roundpg.NewRepository(pgClient)
	roundService := roundsvc.NewService(roundPostgres, packPostgres)

	type roundUseCase struct {
		*roundpg.Repository
		*roundsvc.Service
	}

	roundHandlerV1 := roundv1.NewHandler(&roundUseCase{roundPostgres, roundService}, sessionManager)

	// http
	mux := apptwirp.NewMux([]apptwirp.Handler{
		playerHandlerV1,
		tagHandlerV1,
		authHandlerV1,
		mediaHandlerV1,
		questionHandlerV1,
		packHandlerV1,
		roundHandlerV1,
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
