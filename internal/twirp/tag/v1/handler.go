package v1

import (
	"context"
	"net/http"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/paging"
	"github.com/ysomad/answersuck/internal/pkg/sort"
	"github.com/ysomad/answersuck/internal/twirp"

	pb "github.com/ysomad/answersuck/internal/gen/api/tag/v1"
)

var (
	_ twirp.Handler = &Handler{}
	_ pb.TagService = &Handler{}
)

type UseCase interface {
	Save(context.Context, entity.Tag) error
	GetAll(ctx context.Context, search string, p paging.OffsetParams, s []sort.Sort) (paging.List[entity.Tag], error)
}

type Handler struct {
	tag UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{
		tag: uc,
	}
}

func (h *Handler) Handle(m *http.ServeMux) {
	s := pb.NewTagServiceServer(h)
	m.Handle(s.PathPrefix(), s)
}
