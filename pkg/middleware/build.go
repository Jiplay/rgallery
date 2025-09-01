package middleware

import (
	"net/http"
)

func Build(sha, tag string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("build", sha)
			w.Header().Set("version", tag)
			next.ServeHTTP(w, r)

		})
	}
}
