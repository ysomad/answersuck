package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) LogOut(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	cookie := http.Cookie{
		Name:     session.Cookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
	}

	if err := twirp.SetHTTPResponseHeader(ctx, "Set-Cookie", cookie.String()); err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return new(emptypb.Empty), nil
}
