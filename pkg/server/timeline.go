package server

import (
	"fmt"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/robbymilo/rgallery/pkg/queries"
	"github.com/robbymilo/rgallery/pkg/render"
	"github.com/robbymilo/rgallery/pkg/types"
)

type Media = types.Media
type Days = types.Days
type Filter = types.Filter
type ResponseFilter = types.ResponseFilter

type CacheKey = types.CacheKey

func ServeTimeline(w http.ResponseWriter, r *http.Request) {
	c := r.Context().Value(ConfigKey{}).(Conf)

	params := r.Context().Value(ParamsKey{}).(FilterParams)
	var response ResponseFilter

	if params.Json {

		// cache
		cacheContext := r.Context().Value(CacheKey{}).(map[string]interface{})
		cacheHandle := cacheContext["cache"].(*cache.Cache)
		cacheMap, ok := cacheContext["response"].(ResponseFilter)
		if ok {
			err := render.Render(w, r, cacheMap, "index")
			if err != nil {
				c.Logger.Error("error rendering cached timeline response", "error", err)
			}
			return
		}

		index, total, err := queries.GetTimeline(&params, c)
		if err != nil {
			c.Logger.Error("error getting timeline", "error", err)
		}

		response = ResponseFilter{
			ResponseSegment: index,
			OrderBy:         params.OrderBy,
			Page:            params.Page,
			PageSize:        -1,
			Total:           total,
			Direction:       params.Direction,
			Section:         "timeline",
			Filter: Filter{
				Camera:        params.Camera,
				Lens:          params.Lens,
				Term:          params.Term,
				Mediatype:     params.MediaType,
				Rating:        params.Rating,
				Folder:        params.Folder,
				Subject:       params.Subject,
				Software:      params.Software,
				FocalLength35: params.FocalLength35,
			},
			HideNavFooter: false,
			Meta:          c.Meta,
		}

		err = render.Render(w, r, response, "timeline")
		if err != nil {
			c.Logger.Error("error rendering timeline response", "error", err)
		}

		// setting cache
		cacheHandle.Set(fmt.Sprint(r.URL)+time.Now().Format("2006-01-02"), response, cache.NoExpiration)

	} else {
		response = ResponseFilter{
			Section:       "timeline",
			HideNavFooter: false,
			Meta:          c.Meta,
		}

		err := render.Render(w, r, response, "timeline")
		if err != nil {
			c.Logger.Error("error rendering timeline response", "error", err)
		}
	}
}
