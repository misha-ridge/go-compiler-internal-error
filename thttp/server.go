package thttp

import (
	"context"
	"net"
	"net/http"
)

func Run(start func(spawn func(task func() error)) error) error {
	return nil
}

type Server struct {
}

func NewServer() *Server {
	return nil
}

func (s *Server) Run() error {
	return Run(func(spawn func(task func() error)) error {
		_ = http.Server{
			ConnContext: s.connContext,
		}
		return nil
	})
}

func (s *Server) connContext(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}
