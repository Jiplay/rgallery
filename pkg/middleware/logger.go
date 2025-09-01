package middleware

import (
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/middleware"
)

// log all URLs with request duration
func Logger(c Conf) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			c.Logger.Info("request", "method", r.Method, "url", r.URL.String(),
				"status", ww.Status(), "duration", duration.String())
		})
	}
}
