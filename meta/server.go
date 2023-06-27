package meta

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/misha-ridge/x/parallel"
	"github.com/misha-ridge/x/thttp"
	"github.com/misha-ridge/x/tlog"
	"github.com/misha-ridge/x/tnet"
	"github.com/ridge/must/v2"
)

// IP stores IP used to send requests to meta server
const IP = "169.254.169.254"

// Run runs metadata server
func Run(ctx context.Context, allowedNetwork net.IPNet, addr string, params map[string]string, awsCredentialsSupplier AWSCredentialsSupplier) error {
	l, err := tnet.Listen(addr)
	if err != nil {
		return err
	}

	paramsMarshaled := must.OK1(json.Marshal(params))
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
	router.HandleFunc("/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(paramsMarshaled)
	})
	// FIXME (misha): this works because the AWS client does not really care about
	// the name of the role in the AWS metadata provisioning. But to be exactly like AWS,
	// we may want to pass the role ARN or role name from client and use the exact role name.
	awsCredentialsHandler := newAWSHandler(tlog.Get(ctx), awsCredentialsSupplier, "tectonic", allowedNetwork)
	router.PathPrefix(awsBaseURL).Handler(awsCredentialsHandler)

	return parallel.Run(ctx, func(ctx context.Context, spawn parallel.SpawnFn) error {
		spawn("server", parallel.Fail, thttp.NewServer(l, thttp.StandardMiddleware(router)).Run)
		spawn("iptables", parallel.Fail, func(ctx context.Context) error {
			return ConfigureIPTables(ctx, allowedNetwork, addr)
		})
		return nil
	})
}
