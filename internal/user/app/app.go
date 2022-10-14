package app

import (
	emailv1 "github.com/ysomad/answersuck/internal/user/handler/connect/email/v1"
	"github.com/ysomad/answersuck/internal/user/postgres"
	"github.com/ysomad/answersuck/internal/user/service"
	"github.com/ysomad/answersuck/rpc/user/account/v1/accountv1connect"
	"github.com/ysomad/answersuck/rpc/user/email/v1/emailv1connect"
	"github.com/ysomad/answersuck/rpc/user/password/v1/passwordv1connect"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ysomad/answersuck/internal/user/config"
	accountv1 "github.com/ysomad/answersuck/internal/user/handler/connect/account/v1"
	passwordv1 "github.com/ysomad/answersuck/internal/user/handler/connect/password/v1"
	"github.com/ysomad/answersuck/internal/user/pkg/argon2"

	"github.com/ysomad/answersuck/httpserver"
	"github.com/ysomad/answersuck/logger"
	"github.com/ysomad/answersuck/pgclient"
)

func Run(conf *config.Config) {
	log := logger.New(conf.App.Ver, logger.WithLevel(conf.Log.Level), logger.WithMoscowLocation())

	// dependencies
	postgresClient, err := pgclient.New(conf.PG.URL, pgclient.WithMaxConns(conf.PG.MaxConns))
	if err != nil {
		log.Fatalf("pgclient.New: %w", err)
	}

	passwordHasher := argon2.New()

	// repositories
	accountRepo := postgres.NewAccountRepository(postgresClient)

	// services
	accountService := service.NewAccountService(accountRepo, passwordHasher)

	// handlers
	mux := http.NewServeMux()

	accountV1Server := accountv1.NewServer(log, accountService)
	accountV1Path, accountV1Handler := accountv1connect.NewAccountServiceHandler(accountV1Server)
	mux.Handle(accountV1Path, accountV1Handler)

	passwordV1Server := passwordv1.NewServer(log)
	passwordV1Path, passwordV1Handler := passwordv1connect.NewPasswordServiceHandler(passwordV1Server)
	mux.Handle(passwordV1Path, passwordV1Handler)

	emailV1Server := emailv1.NewServer(log)
	emailV1Path, emailV1Handler := emailv1connect.NewEmailServiceHandler(emailV1Server)
	mux.Handle(emailV1Path, emailV1Handler)

	runHTTPServer(mux, log, conf.HTTP.Port)
}

func runHTTPServer(mux http.Handler, log logger.Logger, port string) {
	log.Infof("starting http server on port %s", port)

	httpServer := httpserver.New(mux, httpserver.WithPort(port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Infof("received signal from httpserver: %s", s.String())
	case err := <-httpServer.Notify():
		log.Infof("got error from http server notify %s", err.Error())
	}

	if err := httpServer.Shutdown(); err != nil {
		log.Infof("got error on http server shutdown %s", err.Error())
	}
}
