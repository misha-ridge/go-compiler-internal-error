package thttp

import (
	"context"
	"net"
	"net/http"

	"github.com/misha-ridge/x/parallel"
	"github.com/misha-ridge/x/time"
	"github.com/misha-ridge/x/tlog"
	"github.com/ridge/must/v2"
	"go.uber.org/zap"
)

const gracefulShutdownTimeout = 5 * time.Second

// Server wraps an HTTP server
type Server struct {
	listener net.Listener
	handler  http.Handler
}

// NewServer creates a Server
func NewServer(listener net.Listener, handler http.Handler) *Server {
	return &Server{
		listener: listener,
		handler:  handler,
	}
}

// Run serves requests until the context is closed, then performs graceful
// shutdown for up to gracefulShutdownTimeout
func (s *Server) Run(ctx context.Context) error {
	return parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		reqCtx, _ := context.WithCancel(ctx) // stays open longer than ctx

		logger := tlog.Get(ctx)

		_ = http.Server{
			Handler:     s.handler,
			ErrorLog:    must.OK1(zap.NewStdLogAt(logger, zap.WarnLevel)),
			BaseContext: func(net.Listener) context.Context { return reqCtx },
			ConnContext: s.connContext,
		}

		return nil
	})
}

// ListenAddr returns the local address of the server's listener
func (s *Server) ListenAddr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) connContext(ctx context.Context, conn net.Conn) context.Context {
	return tlog.With(ctx, zap.Stringer("remoteAddr", conn.RemoteAddr()))
}
