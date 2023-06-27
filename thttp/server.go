package thttp

import (
	"context"
	"net"
	"net/http"
)

func Run(ctx context.Context, start func(spawn func(task func(ctx context.Context) error)) error) error {
	return nil
}

type Server struct {
}

func NewServer() *Server {
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	return Run(ctx, func(spawn func(task func(ctx context.Context) error)) error {
		_ = http.Server{
			ConnContext: s.connContext,
		}
		return nil
	})
}

func (s *Server) connContext(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}
