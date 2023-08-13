package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/paging"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"github.com/ysomad/answersuck/internal/pkg/sort"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	_ apptwirp.Handler = &TagHandler{}
	_ pb.TagService    = &TagHandler{}
)

type TagUseCase interface {
	Save(context.Context, entity.Tag) error
	GetAll(context.Context, paging.Params, []sort.Sort) (paging.List[entity.Tag], error)
}

type TagHandler struct {
	tag     TagUseCase
	session *session.Manager
}

func NewTagHandler(uc TagUseCase, sm *session.Manager) *TagHandler {
	return &TagHandler{
		tag:     uc,
		session: sm,
	}
}

func (h *TagHandler) Handle(m *http.ServeMux) {
	s := pb.NewTagServiceServer(h,
		twirp.WithServerHooks(hooks.WithSession(h.session)))
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}

func (h *TagHandler) CreateTag(
	ctx context.Context,
	r *pb.CreateTagRequest) (*pb.CreateTagResponse, error) {
	session, err := common.CheckPlayerVerification(ctx)
	if err != nil {
		return nil, err
	}

	if r.TagName == "" {
		return nil, twirp.RequiredArgumentError("tag_name")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	tag := entity.Tag{
		Name:       r.TagName,
		Author:     session.User.ID,
		CreateTime: time.Now(),
	}

	if err := h.tag.Save(ctx, tag); err != nil {
		if errors.Is(err, apperr.TagAlreadyExists) {
			return nil, twirp.AlreadyExists.Error(apperr.MsgTagAlreadyExists)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreateTagResponse{
		Tag: &pb.Tag{
			Name:       tag.Name,
			Author:     tag.Author,
			CreateTime: timestamppb.New(tag.CreateTime),
		},
	}, nil
}

func (h *TagHandler) ListTags(
	ctx context.Context, r *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	sorts, err := sort.NewSortList(r.OrderBy)
	if err != nil {
		return nil, twirp.InvalidArgumentError("order_by", err.Error())
	}

	tagList, err := h.tag.GetAll(ctx, paging.Params{
		PageSize:  r.PageSize,
		PageToken: r.PageToken,
	}, sorts)
	if err != nil {
		if errors.Is(err, paging.ErrInvalidToken) {
			return nil, twirp.InvalidArgumentError("page_token", err.Error())
		}

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
		Tags:          tags,
		NextPageToken: tagList.NextPageToken,
	}, nil
}
