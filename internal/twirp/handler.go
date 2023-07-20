package twirp

import "net/http"

type Handler interface {
	Handle(*http.ServeMux)
}

// NewMux creates new mux and attaches handlers to it.
func NewMux(handlers []Handler) *http.ServeMux {
	m := http.NewServeMux()
	for _, h := range handlers {
		h.Handle(m)
	}
	return m
}
