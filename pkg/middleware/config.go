package middleware

import (
	"context"
	"net/http"

	"github.com/robbymilo/rgallery/pkg/types"
)

type ConfigKey = types.ConfigKey

func Config(config Conf) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := context.WithValue(r.Context(), ConfigKey{}, config)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
