package v1

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type packageHandler struct {
	log *zap.Logger
}

func newPackageMux(d *Deps) *chi.Mux {
	// h := packageHandler{}
	m := chi.NewMux()

	return m
}
