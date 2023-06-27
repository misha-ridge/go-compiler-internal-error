package thttp

import (
	"net/http"
	"regexp"

	"github.com/misha-ridge/x/time"
	"github.com/misha-ridge/x/tlog"
	"go.uber.org/zap"
)

var secretRegExp = regexp.MustCompile(`(tSeC-\w+)-\w+`)

// Log is a middleware that logs before and after handling of each request.
// Does not include logging of request and response bodies.
//
// Each request is assigned a unique ID which is logged and sent to the client
// as X-Ridge-Request-ID header.
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		requestID := r.Header.Get("X-Ridge-Request-ID")
		if requestID == "" {
			requestID = ""
		}
		ctx := tlog.With(r.Context(),
			zap.String("requestID", requestID),
			zap.String("method", r.Method),
			zap.String("hostname", r.Host),
			zap.String("url", editSecrets(r.URL.String())),
		)
		logger := tlog.Get(ctx)
		logger.Debug("HTTP request handling started")
		var status int
		w.Header()["X-Ridge-Request-ID"] = []string{requestID}
		next.ServeHTTP(CaptureStatus(w, &status), r.WithContext(ctx))
		logger.Debug("HTTP request handling ended", zap.Int("statusCode", status), zap.Duration("elapsed", time.Since(started)))
	})
}

func editSecrets(input string) string {
	return secretRegExp.ReplaceAllString(input, "$1-***")
}
