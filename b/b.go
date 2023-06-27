package b

import (
	"context"
	"net"
	"net/http"
)

func R(func(func(func())) error) error {
	return nil
}

type S struct {
}

func NewS() *S {
	return nil
}

func (s *S) R() {
	R(func(func(func())) error {
		_ = http.Server{
			ConnContext: s.connContext,
		}
		return nil
	})
}

func (s *S) connContext(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}
