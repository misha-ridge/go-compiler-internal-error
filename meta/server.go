package meta

import (
	"context"
	"net"

	"github.com/misha-ridge/x/thttp"
)

// Run runs metadata server
func Run(ctx context.Context, allowedNetwork net.IPNet, addr string, params map[string]string) error {
	return thttp.Run(ctx, func(ctx context.Context, spawn thttp.SpawnFn) error {
		spawn("server", thttp.Fail, thttp.NewServer(nil, nil).Run)
		return nil
	})
}
