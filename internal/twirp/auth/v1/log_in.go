package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/auth/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	SessionCookieName = "sid"
	HeaderSetCookie   = "Set-Cookie"
)

func (h *Handler) LogIn(ctx context.Context, p *pb.LogInRequest) (*emptypb.Empty, error) {
	cookie := http.Cookie{
		Name:     SessionCookieName,
		Value:    "test",
		Expires:  time.Now().Add(time.Hour),
		Secure:   false,
		HttpOnly: false,
	}

	twirp.SetHTTPResponseHeader(ctx, HeaderSetCookie, cookie.String())

	return new(emptypb.Empty), nil
}
