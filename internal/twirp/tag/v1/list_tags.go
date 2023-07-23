package v1

import (
	"context"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/tag/v1"
	"github.com/ysomad/answersuck/internal/pkg/paging"
	"github.com/ysomad/answersuck/internal/pkg/sort"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) ListTags(ctx context.Context, r *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}

	sorts, err := sort.NewSortList(r.OrderBy)
	if err != nil {
		return nil, twirp.InvalidArgumentError("order_by", err.Error())
	}

	tagList, err := h.tag.GetAll(ctx, paging.OffsetParams{
		Offset: r.Offset,
		Limit:  r.Limit,
	}, sorts)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	tags := make([]*pb.Tag, len(tagList.Items))

	for i, t := range tagList.Items {
		tags[i] = &pb.Tag{
			Name:       t.Name,
			Author:     t.Author,
			CreateTime: timestamppb.New(t.CreateTime),
		}
	}

	return &pb.ListTagsResponse{
		Tags:    tags,
		HasNext: tagList.HasNext,
	}, nil
}
