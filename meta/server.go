package meta

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/misha-ridge/x/parallel"
	"github.com/misha-ridge/x/thttp"
)

// IP stores IP used to send requests to meta server
const IP = "169.254.169.254"

// Run runs metadata server
func Run(ctx context.Context, allowedNetwork net.IPNet, addr string, params map[string]string) error {
	var l net.Listener

	router := mux.NewRouter()
	router.HandleFunc("/v1/{paramKey}", func(w http.ResponseWriter, r *http.Request) {
		paramKey := strings.ReplaceAll(mux.Vars(r)["paramKey"], "-", "_")
		if val, ok := params[paramKey]; ok {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte(val))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	return parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		spawn("server", parallel.Fail, thttp.NewServer(l, thttp.StandardMiddleware(router)).Run)
		return nil
	})
}
