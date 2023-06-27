package meta

import (
	"context"
	"net"
	"net/http"

	"github.com/misha-ridge/x/parallel"
	"github.com/misha-ridge/x/thttp"
)

// Run runs metadata server
func Run(ctx context.Context, allowedNetwork net.IPNet, addr string, params map[string]string) error {
	var l net.Listener
	var router http.Handler

	x := thttp.NewServer(l, router).Run

	return parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		spawn("server", parallel.Fail, x)
		return nil
	})
}
