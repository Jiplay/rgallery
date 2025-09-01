package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/robbymilo/rgallery/pkg/types"
)

type CacheKey = types.CacheKey

func Cache(cache *cache.Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx context.Context

			var cacheMap = make(map[string]interface{})
			cacheMap["cache"] = cache

			response, found := cache.Get(fmt.Sprint(r.URL) + time.Now().Format("2006-01-02"))
			if found {
				w.Header().Set("Cache-Status", "HIT")
				cacheMap["response"] = response

			} else {
				w.Header().Set("Cache-Status", "MISS")

			}

			ctx = context.WithValue(r.Context(), CacheKey{}, cacheMap)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
