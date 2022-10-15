package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ysomad/answersuck/internal/peasant/config"
	"github.com/ysomad/answersuck/internal/peasant/pkg/argon2"
	"github.com/ysomad/answersuck/internal/peasant/postgres"
	"github.com/ysomad/answersuck/internal/peasant/service"

	"github.com/ysomad/answersuck/rpc/peasant/v1/v1connect"

	accountv1 "github.com/ysomad/answersuck/internal/peasant/handler/connect/account/v1"
	emailv1 "github.com/ysomad/answersuck/internal/peasant/handler/connect/email/v1"
	passwordv1 "github.com/ysomad/answersuck/internal/peasant/handler/connect/password/v1"

	"github.com/ysomad/answersuck/httpserver"
	"github.com/ysomad/answersuck/logger"
	"github.com/ysomad/answersuck/pgclient"
)

func Run(conf *config.Config) {
	log := logger.New(
		conf.App.Ver,
		logger.WithLevel(conf.Log.Level),
		logger.WithLocation(conf.Log.TimeLoc),
		logger.WithSkipFrameCount(1),
	)

	// dependencies
	pg, err := pgclient.New(conf.PG.URL, pgclient.WithMaxConns(conf.PG.MaxConns))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer pg.Close()

	passwordHasher := argon2.New()

	// repositories
	accountRepo := postgres.NewAccountRepository(pg)

	// services
	accountService, err := service.NewAccountService(accountRepo, passwordHasher)
	if err != nil {
		log.Fatalf("service.NewAccountService: %s", err.Error())
	}

	emailService, err := service.NewEmailService(accountRepo, passwordHasher)
	if err != nil {
		log.Fatalf("service.NewEmailService: %s", err.Error())
	}

	// http
	mux := http.NewServeMux()

	accountV1Server := accountv1.NewServer(log, accountService)
	accountV1Path, accountV1Handler := v1connect.NewAccountServiceHandler(accountV1Server)
	mux.Handle(accountV1Path, accountV1Handler)

	passwordV1Server := passwordv1.NewServer(log)
	passwordV1Path, passwordV1Handler := v1connect.NewPasswordServiceHandler(passwordV1Server)
	mux.Handle(passwordV1Path, passwordV1Handler)

	emailV1Server := emailv1.NewServer(log, emailService)
	emailV1Path, emailV1Handler := v1connect.NewEmailServiceHandler(emailV1Server)
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
