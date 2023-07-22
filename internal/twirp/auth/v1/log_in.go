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

func (h *Handler) LogIn(ctx context.Context, p *pb.LogInRequest) (*emptypb.Empty, error) {
	if err := h.validateLogInRequest(p); err != nil {
		return nil, err
	}

	sid := appctx.GetSessionID(ctx)
	if sid != "" {
		return new(emptypb.Empty), nil
	}

	s, err := h.auth.LogIn(ctx, p.Login, p.Password, appctx.GetFootPrint(ctx))
	if err != nil {
		if errors.Is(err, apperr.ErrPlayerNotFound) || errors.Is(err, apperr.ErrNotAuthorized) {
			return nil, twirp.Unauthenticated.Error(apperr.ErrNotAuthorized.Error())
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

func (h *Handler) validateLogInRequest(p *pb.LogInRequest) error {
	if p.Login == "" {
		return twirp.RequiredArgumentError("login")
	}

	if p.Password == "" {
		return twirp.RequiredArgumentError("password")
	}

	return nil
}
