package server

import (
	"net/http"

	cache "github.com/patrickmn/go-cache"
	"github.com/robbymilo/rgallery/pkg/middleware"
)

func Purge(w http.ResponseWriter, r *http.Request, cache *cache.Cache, c Conf) {

	cache.Flush()
	middleware.RemoveEtags()

	c.Logger.Info("cache purged")

	_, err := w.Write([]byte("ok"))
	if err != nil {
		c.Logger.Error("error writing purge status", "error", err)
	}
}
