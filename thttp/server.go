package thttp

import (
	"context"
	"net"
	"net/http"
)

type SpawnFn func(task func(ctx context.Context) error)

func Run(ctx context.Context, start func(spawn SpawnFn) error) error {
	return nil
}

type Server struct {
}

func NewServer(listener net.Listener, handler http.Handler) *Server {
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	return Run(ctx, func(spawn SpawnFn) error {
		_ = http.Server{
			ConnContext: s.connContext,
		}
		return nil
	})
}

func (s *Server) connContext(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}
