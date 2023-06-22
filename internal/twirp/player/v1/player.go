package playerv1

import (
	"context"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/ysomad/answersuck/internal/gen/proto/player/v1"
)

var _ pb.PlayerService = &playerHandler{}

type playerHandler struct{}

func NewPlayerHandler() *playerHandler {
	return &playerHandler{}
}

func (h *playerHandler) CreatePlayer(ctx context.Context, r *pb.CreatePlayerRequest) (*emptypb.Empty, error) {
	if err := r.Validate(); err != nil {
		if pberr, ok := err.(pb.CreatePlayerRequestValidationError); ok {
			return nil, twirp.InvalidArgumentError(pberr.Field(), pberr.Error())
		}

		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}

	return new(emptypb.Empty), twirp.InternalError("not implemented")
}

func (h *playerHandler) GetPlayer(ctx context.Context, r *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	c := http.Cookie{
		Name:     "sid",
		Value:    "test",
		Path:     "http://localhost:8080",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Hour),
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	}

	if err := twirp.SetHTTPResponseHeader(ctx, "Set-Cookie", c.String()); err != nil {
		return nil, twirp.InternalError("cookie not set")
	}

	return nil, twirp.InternalError("not implemented")
}
