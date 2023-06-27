package thttp

import (
	"context"
	"net"
	"net/http"

	"github.com/misha-ridge/x/parallel"
)

// Server wraps an HTTP server
type Server struct {
}

// NewServer creates a Server
func NewServer(listener net.Listener, handler http.Handler) *Server {
	return &Server{}
}

// Run serves requests until the context is closed, then performs graceful
// shutdown for up to gracefulShutdownTimeout
func (s *Server) Run(ctx context.Context) error {
	return parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		_ = http.Server{
			ConnContext: s.connContext,
		}

		return nil
	})
}

func (s *Server) connContext(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}
