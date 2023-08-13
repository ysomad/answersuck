package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/session"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
)

var (
	_ apptwirp.Handler = &PackHandler{}
	_ pb.PackService   = &PackHandler{}
)

type PackUseCase interface {
	Save(ctx context.Context, p *entity.Pack, tags []string) (packID int32, err error)
	GetWithTags(ctx context.Context, packID int32) (*entity.PackWithTags, error)
}

type PackHandler struct {
	pack    PackUseCase
	session *session.Manager
}

func NewPackHandler(uc PackUseCase, sm *session.Manager) *PackHandler {
	return &PackHandler{
		pack:    uc,
		session: sm,
	}
}

func (h *PackHandler) Handle(m *http.ServeMux) {
	s := pb.NewPackServiceServer(h,
		twirp.WithServerHooks(hooks.WithSession(h.session)))
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}

func (h *PackHandler) CreatePack(
	ctx context.Context,
	r *pb.CreatePackRequest) (*pb.CreatePackResponse, error) {
	session, err := common.CheckPlayerVerification(ctx)
	if err != nil {
		return nil, err
	}

	if r.PackName == "" {
		return nil, twirp.RequiredArgumentError("pack_name")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	packID, err := h.pack.Save(ctx, &entity.Pack{
		Name:       r.PackName,
		Author:     session.User.ID,
		CoverURL:   r.CoverUrl,
		CreateTime: time.Now(),
	}, r.Tags)
	if err != nil {
		if errors.Is(err, apperr.MediaNotFound) {
			return nil, twirp.InvalidArgument.Error(apperr.MsgPackCoverNotFound)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreatePackResponse{PackId: packID}, nil
}

func (h *PackHandler) GetPack(
	ctx context.Context, r *pb.GetPackRequest) (*pb.GetPackResponse, error) {
	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	p, err := h.pack.GetWithTags(ctx, r.PackId)
	if err != nil {
		if errors.Is(err, apperr.PackNotFound) {
			return nil, twirp.NotFoundError(apperr.MsgPackNotFound)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.GetPackResponse{
		Pack: &pb.Pack{
			Id:          p.ID,
			Name:        p.Name,
			Author:      p.Author,
			IsPublished: p.Published,
			CoverUrl:    p.CoverURL,
			CreateTime:  timestamppb.New(p.CreateTime),
		},
		Tags: p.Tags,
	}, nil
}

func (h *PackHandler) PublishPack(
	ctx context.Context,
	r *pb.PublishPackRequest) (*pb.PublishPackResponse, error) {
	return nil, nil
}
