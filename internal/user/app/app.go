package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ysomad/answersuck/internal/user/config"
	"github.com/ysomad/go-packages/httpserver"
)

func Run(conf *config.Config) {
	accountServer := newAccountServiceServer()
	emailServer := newEmailServiceServer()
	passwordServer := newPasswordServiceServer()

	mux := http.NewServeMux()

	mux.Handle(accountServer.PathPrefix(), accountServer)
	mux.Handle(emailServer.PathPrefix(), emailServer)
	mux.Handle(passwordServer.PathPrefix(), passwordServer)

	runHTTPServer(mux, &conf.HTTP)
}

func runHTTPServer(mux http.Handler, conf *config.HTTP) {
	log.Printf("running http server at port %s", conf.Port)

	httpServer := httpserver.New(mux, httpserver.WithPort(conf.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		log.Printf("app - Run - httpServer.Notify: %s", err.Error())
	}

	if err := httpServer.Shutdown(); err != nil {
		log.Printf("app - Run - httpServer.Shutdown: %s", err.Error())
	}
}
