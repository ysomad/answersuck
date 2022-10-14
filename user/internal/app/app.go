package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ysomad/answersuck/user/internal/config"
	"github.com/ysomad/answersuck/user/internal/handler/connect/account"
	"github.com/ysomad/answersuck/user/internal/repository"
	"github.com/ysomad/answersuck/user/internal/service"

	"github.com/ysomad/answersuck/user/internal/gen/proto/account/accountconnect"

	"github.com/ysomad/answersuck/pkg/argon2"
	"github.com/ysomad/answersuck/pkg/httpserver"
	"github.com/ysomad/answersuck/pkg/logger"
	"github.com/ysomad/answersuck/pkg/pgclient"
)

func Run(conf *config.Config) {
	log := logger.New(conf.App.Ver, logger.WithLevel(conf.Log.Level), logger.WithMoscowLocation())

	mux := http.NewServeMux()

	postgresClient, err := pgclient.New(conf.PG)
	if err != nil {
		log.Fatalf("pgclient.New: %w", err)
	}

	passwordHasher := argon2.New()

	accountRepository := repository.NewAccountRepo()
	accountService := service.NewAccountService(accountRepository, passwordHasher)
	accountServer := account.NewServer(log, accountService)
	accountPath, accountHandler := accountconnect.NewAccountServiceHandler(accountServer)

	mux.Handle(accountPath, accountHandler)

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
