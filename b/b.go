package b

import (
	"context"
	"net"
	"net/http"
)

func R(func(func(func() error)) error) error {
	return nil
}

type S struct {
}

func NewS() *S {
	return nil
}

func (s *S) R() error {
	return R(func(func(func() error)) error {
		_ = http.Server{
			ConnContext: s.connContext,
		}
		return nil
	})
}

func (s *S) connContext(ctx context.Context, conn net.Conn) context.Context {
	return ctx
}
