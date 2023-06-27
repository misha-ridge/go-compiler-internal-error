package thttp

import (
	"context"
	"net/http"

	"github.com/misha-ridge/x/parallel"
)

// Recover is a middleware that catches and logs panics from HTTP handlers
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := parallel.RunTask(r.Context(), func(ctx context.Context) error {
			next.ServeHTTP(w, r)
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			select {
			case r.Context().Value(panicKey).(chan error) <- err:
			default:
			}
		}
	})
}
