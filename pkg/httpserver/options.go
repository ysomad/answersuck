package httpserver

import (
	"net"
	"time"
)

type Option func(*server)

func Port(port string) Option {
	return func(s *server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.shutdownTimeout = timeout
	}
}
