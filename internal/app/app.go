package app

import (
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/ysomad/answersuck/internal/config"

	mediapg "github.com/ysomad/answersuck/internal/postgres/media"
	packpg "github.com/ysomad/answersuck/internal/postgres/pack"
	playerpg "github.com/ysomad/answersuck/internal/postgres/player"
	questionpg "github.com/ysomad/answersuck/internal/postgres/question"
	roundpg "github.com/ysomad/answersuck/internal/postgres/round"
	"github.com/ysomad/answersuck/internal/postgres/roundquestion"
	roundtopicpg "github.com/ysomad/answersuck/internal/postgres/roundtopic"
	tagpg "github.com/ysomad/answersuck/internal/postgres/tag"
	topicpg "github.com/ysomad/answersuck/internal/postgres/topic"

	authsvc "github.com/ysomad/answersuck/internal/service/auth"
	"github.com/ysomad/answersuck/internal/service/pack"
	playersvc "github.com/ysomad/answersuck/internal/service/player"
	roundsvc "github.com/ysomad/answersuck/internal/service/round"

	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	authv1 "github.com/ysomad/answersuck/internal/twirp/auth/v1"
	editorv1 "github.com/ysomad/answersuck/internal/twirp/editor/v1"
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
	tagHandlerV1 := editorv1.NewTagHandler(tagPostgres, sessionManager)

	// auth
	authService := authsvc.NewService(sessionManager, playerService)
	authHandlerV1 := authv1.NewAuthHandler(authService)

	// media
	mediaPostgres := mediapg.NewRepository(pgClient)
	mediaHandlerV1 := editorv1.NewMediaHandler(mediaPostgres, sessionManager)

	// question
	questionPostgres := questionpg.NewRepository(pgClient)
	questionHandlerV1 := editorv1.NewQuestionHandler(questionPostgres, sessionManager)

	// pack
	packPostgres := packpg.NewRepository(pgClient)
	packSvc := pack.NewService(packPostgres)
	packHandlerV1 := editorv1.NewPackHandler(packPostgres, sessionManager)

	// topic
	topicPostgres := topicpg.NewRepository(pgClient)
	topicHandlerV1 := editorv1.NewTopicHandler(topicPostgres, sessionManager)

	// roundTopic
	roundTopicPostgres := roundtopicpg.NewRepository(pgClient)

	// round
	roundPostgres := roundpg.NewRepository(pgClient)
	roundService := roundsvc.NewService(roundPostgres, packSvc, roundTopicPostgres)

	type roundUseCase struct {
		*roundpg.Repository
		*roundsvc.Service
	}

	roundHandlerV1 := editorv1.NewRoundHandler(&roundUseCase{roundPostgres, roundService}, sessionManager)

	// round question
	roundQuestionPostgres := roundquestion.NewRepository(pgClient)
	roundQuestionHandlerV1 := editorv1.NewRoundQuestionHandler(roundQuestionPostgres, sessionManager)

	// http
	mux := apptwirp.NewMux([]apptwirp.Handler{
		playerHandlerV1,
		tagHandlerV1,
		authHandlerV1,
		mediaHandlerV1,
		questionHandlerV1,
		packHandlerV1,
		roundHandlerV1,
		topicHandlerV1,
		roundQuestionHandlerV1,
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
