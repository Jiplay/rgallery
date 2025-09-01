package rgallery

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/robbymilo/rgallery/pkg/dist"
	"github.com/robbymilo/rgallery/pkg/fonts"
	"github.com/robbymilo/rgallery/pkg/render"
	"github.com/robbymilo/rgallery/pkg/static"
)

func Dist(h http.Handler, c Conf) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		content, err := dist.DistDir.ReadFile(strings.TrimPrefix(r.URL.Path, "/dist/"))
		if err != nil {
			c.Logger.Error("error reading file", "error", err)
		}

		e := render.GenerateEtag(string(content))

		w.Header().Set("Etag", fmt.Sprintf("\"%s\"", e))
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Static(h http.Handler, c Conf) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		content, err := static.StaticDir.ReadFile(strings.TrimPrefix(r.URL.Path, "/static/"))
		if err != nil {
			c.Logger.Error("error reading file", "error", err)
		}

		e := render.GenerateEtag(fmt.Sprint(content))
		w.Header().Set("Etag", fmt.Sprintf("\"%s\"", e))
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Fonts(h http.Handler, c Conf) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		content, err := fonts.FontDir.ReadFile(strings.TrimPrefix(r.URL.Path, "/fonts/"))
		if err != nil {
			c.Logger.Error("error reading file", "error", err)
		}

		e := render.GenerateEtag(fmt.Sprint(content))
		w.Header().Set("Etag", fmt.Sprintf("\"%s\"", e))
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
