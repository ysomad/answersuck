package httpserver

import (
	"net"
	"time"
)

type Option func(*server)

func WithPort(port string) Option {
	return func(s *server) {
		s.srv.Addr = net.JoinHostPort("", port)
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.srv.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.srv.WriteTimeout = timeout
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.shutdownTimeout = timeout
	}
}
