package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/robbymilo/rgallery/pkg/queries"
	"github.com/robbymilo/rgallery/pkg/render"
	"github.com/robbymilo/rgallery/pkg/types"
)

type ResponseGear = types.ResponseGear
type Conf = types.Conf
type PrevNext = types.PrevNext

func ServeMedia(w http.ResponseWriter, r *http.Request) {
	c := r.Context().Value(ConfigKey{}).(Conf)
	params := r.Context().Value(ParamsKey{}).(FilterParams)

	h, err := DecodeURL(chi.URLParam(r, "hash"))
	if err != nil {
		c.Logger.Error("error decoding hash", "error", err)
	}
	hash := GetHash(h)
	collection := chi.URLParam(r, "collection")
	slug, err := DecodeURL(getAfter5thSlash(r.URL.Path))
	if err != nil {
		c.Logger.Error("error decoding slug", "error", err)
	}
	if slug == "root" {
		slug = "."
	}
	media, err := queries.GetSingleMediaItem(hash, c)
	if err != nil {
		c.Logger.Error("error getting single media item", "error", err)
	}

	if media.Path == "" {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("404\n"))
		if err != nil {
			c.Logger.Error("error writing media 404 response", "error", err)

		}
		return
	}

	var column string
	if collection == "tag" {
		column = "tag"
	} else if collection == "folder" {
		column = "folder"
	}

	rating := 0
	if collection == "favorites" {
		rating = 5
	} else if params.Rating > 0 {
		rating = params.Rating
	}

	previous, err := queries.GetPrevious(media.Date, hash, column, slug, rating, params, c)
	if err != nil {
		c.Logger.Error("error getting previous media items", "error", err)
	}
	total_next := 6 - len(previous)
	next, err := queries.GetNext(media.Date, hash, column, slug, rating, total_next, params, previous, c)
	if err != nil {
		c.Logger.Error("error getting next media items", "error", err)
	}

	response := ResponseMedia{
		Media:         media,
		Previous:      previous,
		Next:          next,
		Collection:    collection,
		Slug:          slug,
		Section:       "media",
		HideNavFooter: false,
		TileServer:    c.TileServer,
		Meta:          c.Meta,
	}

	err = render.Render(w, r, response, "media")
	if err != nil {
		c.Logger.Error("error rendering media response", "error", err)
	}

}

func getAfter5thSlash(s string) string {
	index := -1
	slashCount := 0

	for i, c := range s {
		if c == '/' {
			slashCount++
			if slashCount == 5 {
				index = i
				break
			}
		}
	}

	if index == -1 || index+1 >= len(s) {
		return "" // Less than 5 slashes or nothing after
	}

	return s[index+1:]
}
