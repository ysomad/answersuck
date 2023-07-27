package pack

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetPack(ctx context.Context, r *pb.GetPackRequest) (*pb.GetPackResponse, error) {
	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	p, err := h.pack.GetOne(ctx, r.PackId)
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
	}, nil
}
