package app

import (
	"github.com/ysomad/answersuck/internal/user/handler"
	"github.com/ysomad/answersuck/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ysomad/answersuck/internal/user/config"
	"github.com/ysomad/go-packages/httpserver"
)

func Run(conf *config.Config) {
	log := logger.New(conf.App.Ver, logger.WithLevel(conf.Log.Level), logger.WithMoscowLocation())

	twirpMux := handler.NewTwirpMux(log, &conf.Twirp)

	runHTTPServer(twirpMux, log, conf.HTTP.Port)
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
