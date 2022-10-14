package handler

import (
	"net/http"

	"github.com/twitchtv/twirp"

	pb "github.com/ysomad/answersuck/rpc/user"

	"github.com/ysomad/answersuck/internal/user/config"
	"github.com/ysomad/answersuck/internal/user/handler/twirp/account"
	"github.com/ysomad/answersuck/internal/user/handler/twirp/email"
	"github.com/ysomad/answersuck/internal/user/handler/twirp/password"
	"github.com/ysomad/answersuck/internal/user/postgres"
	"github.com/ysomad/answersuck/internal/user/service"

	"github.com/ysomad/answersuck/pkg/argon2"
	"github.com/ysomad/answersuck/pkg/logger"
)

func NewTwirpMux(log logger.Logger, conf *config.Twirp) *http.ServeMux {
	handlerPrefix := twirp.WithServerPathPrefix(conf.Prefix)

	passwordHasher := argon2.New()

	accountRepository := postgres.NewAccountRepository()
	accountService := service.NewAccount(accountRepository, passwordHasher)
	accountServer := account.NewServer(log, accountService)
	accountHandler := pb.NewAccountServiceServer(accountServer, handlerPrefix)

	emailServer := email.NewServer()
	emailHandler := pb.NewEmailServiceServer(emailServer, handlerPrefix)

	passwordServer := password.NewServer()
	passwordHandler := pb.NewPasswordServiceServer(passwordServer, handlerPrefix)

	mux := http.NewServeMux()

	mux.Handle(accountHandler.PathPrefix(), accountHandler)
	mux.Handle(emailHandler.PathPrefix(), emailHandler)
	mux.Handle(passwordHandler.PathPrefix(), passwordHandler)

	return mux
}
