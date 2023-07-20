package v1

import (
	"net/http"

	pb "github.com/ysomad/answersuck/internal/gen/api/auth/v1"
	"github.com/ysomad/answersuck/internal/twirp"
)

var (
	_ twirp.Handler  = &Handler{}
	_ pb.AuthService = &Handler{}
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(m *http.ServeMux) {
	s := pb.NewAuthServiceServer(h)
	m.Handle(s.PathPrefix(), s)
}
