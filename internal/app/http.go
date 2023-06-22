package app

import (
	"net/http"

	"github.com/twitchtv/twirp"

	pbv1 "github.com/ysomad/answersuck/internal/gen/proto/player/v1"
	playerv1 "github.com/ysomad/answersuck/internal/twirp/player/v1"
)

func newServeMux(prefix string) *http.ServeMux {
	srvPrefix := twirp.WithServerPathPrefix(prefix)

	playerHandlerV1 := playerv1.NewPlayerHandler()
	playerServerV1 := pbv1.NewPlayerServiceServer(playerHandlerV1, srvPrefix)

	emailHandlerV1 := playerv1.NewEmailHandler()
	emailServerV1 := pbv1.NewEmailServiceServer(emailHandlerV1, srvPrefix)

	passwordHandlerV1 := playerv1.NewPasswordHandler()
	passwordServerV1 := pbv1.NewPasswordServiceServer(passwordHandlerV1, srvPrefix)

	mux := http.NewServeMux()

	// player v1
	mux.Handle(playerServerV1.PathPrefix(), playerServerV1)
	mux.Handle(emailServerV1.PathPrefix(), emailServerV1)
	mux.Handle(passwordServerV1.PathPrefix(), passwordServerV1)

	return mux
}
