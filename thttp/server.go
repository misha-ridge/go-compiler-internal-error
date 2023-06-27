package thttp

import (
	"context"
	"net"
	"net/http"
)

type Group struct {
	done chan struct{}
}

func NewGroup(ctx context.Context) *Group {
	g := new(Group)
	g.done = make(chan struct{})
	return g
}

func RunTask(ctx context.Context, task func(ctx context.Context) error) (err error) {
	return task(ctx)
}

func (g *Group) runTask(task func(ctx context.Context) error) {
	_ = RunTask(nil, task)
}

func (g *Group) exit(err error) {
}

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
