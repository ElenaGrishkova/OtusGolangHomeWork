package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			StatusCode:     200,
		}
		next.ServeHTTP(recorder, r)
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		logger.Info(fmt.Sprintf(
			"Request to %s %s %s %s %d %s \"%s\"",
			ip,
			r.Method,
			r.RequestURI,
			r.Proto,
			recorder.StatusCode,
			time.Since(start),
			r.UserAgent(),
		))
	})
}
