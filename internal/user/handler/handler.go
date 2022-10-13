package handler

import (
	"net/http"

	"github.com/twitchtv/twirp"

	pb "github.com/ysomad/answersuck/rpc/user"

	"github.com/ysomad/answersuck/internal/user/config"
	"github.com/ysomad/answersuck/internal/user/handler/twirp/account"
	"github.com/ysomad/answersuck/internal/user/handler/twirp/email"
	"github.com/ysomad/answersuck/internal/user/handler/twirp/password"
)

func NewTwirpMux(conf *config.Twirp) *http.ServeMux {
	handlerPrefix := twirp.WithServerPathPrefix(conf.Prefix)

	accountServer := account.NewServer()
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
