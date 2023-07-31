package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/auth/v1"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) LogIn(ctx context.Context, r *pb.LogInRequest) (*emptypb.Empty, error) {
	if r.Login == "" {
		return nil, twirp.RequiredArgumentError("login")
	}

	if r.Password == "" {
		return nil, twirp.RequiredArgumentError("password")
	}

	if sid, ok := appctx.GetSessionID(ctx); ok && sid != "" {
		return new(emptypb.Empty), nil
	}

	fp, ok := appctx.GetFootPrint(ctx)
	if !ok {
		return nil, twirp.InvalidArgument.Error("footprint not found in context")
	}

	s, err := h.auth.LogIn(ctx, r.Login, r.Password, fp)
	if err != nil {
		if errors.Is(err, apperr.PlayerNotFound) || errors.Is(err, apperr.InvalidCredentials) {
			return nil, twirp.Unauthenticated.Error(apperr.InvalidCredentials.Error())
		}

		return nil, twirp.InternalError(err.Error())
	}

	cookie := http.Cookie{
		Name:     session.Cookie,
		Value:    s.ID,
		Path:     "/",
		Expires:  s.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
	}

	if err = twirp.SetHTTPResponseHeader(ctx, "Set-Cookie", cookie.String()); err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return new(emptypb.Empty), nil
}