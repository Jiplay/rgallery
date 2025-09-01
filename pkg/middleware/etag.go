package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var etags = make(map[string]string)

func PersistEtag(key, etag string) {
	var now = time.Now().UTC().Format("20060102")
	etags[key+now] = etag
}

func GetPersistedEtag(key string) string {
	var now = time.Now().UTC().Format("20060102")
	return etags[key+now]
}

func RemoveEtags() {
	etags = make(map[string]string)
}

func Etag(c Conf) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ConfigKey{}, c)
			var user UserKey
			if r.Context().Value(UserKey{}) != nil {
				user = r.Context().Value(UserKey{}).(UserKey)
			}

			params := r.Context().Value(ParamsKey{}).(FilterParams)
			etag := GetPersistedEtag(fmt.Sprint(r.URL) + fmt.Sprint(user) + fmt.Sprint(params.Json))
			if etag != "" && fmt.Sprintf("\"%s\"", etag) == r.Header.Get("If-None-Match") && !c.Dev {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
