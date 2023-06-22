package app

import (
	"net/http"

	"github.com/twitchtv/twirp"

	pbv1 "github.com/ysomad/answersuck/internal/gen/proto/player/v1"
)

func newServeMux(prefix string, c handlerContainer) *http.ServeMux {
	srvPrefix := twirp.WithServerPathPrefix(prefix)

	playerServerV1 := pbv1.NewPlayerServiceServer(c.playerV1, srvPrefix)
	emailServerV1 := pbv1.NewEmailServiceServer(c.emailV1, srvPrefix)
	passwordServerV1 := pbv1.NewPasswordServiceServer(c.passwordV1, srvPrefix)

	mux := http.NewServeMux()

	// player v1
	mux.Handle(playerServerV1.PathPrefix(), playerServerV1)
	mux.Handle(emailServerV1.PathPrefix(), emailServerV1)
	mux.Handle(passwordServerV1.PathPrefix(), passwordServerV1)

	return mux
}
