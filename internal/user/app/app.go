package app

import (
	"github.com/ysomad/answersuck/internal/user/config"
	"github.com/ysomad/answersuck/internal/user/twirp/account"
	pb "github.com/ysomad/answersuck/rpc/user"
	"log"
	"net/http"
)

func Run(conf *config.Config) {
	mux := http.NewServeMux()

	accountSrv := pb.NewAccountServiceServer(account.NewServer())
	mux.Handle(accountSrv.PathPrefix(), accountSrv)

	log.Fatal(http.ListenAndServe(":"+conf.HTTP.Port, mux))
}
