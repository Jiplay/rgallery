package rgallery

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	chiprometheus "github.com/766b/chi-prometheus"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	cache "github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robbymilo/rgallery/pkg/config"
	"github.com/robbymilo/rgallery/pkg/dist"
	"github.com/robbymilo/rgallery/pkg/fonts"
	"github.com/robbymilo/rgallery/pkg/metrics"
	"github.com/robbymilo/rgallery/pkg/middleware"
	"github.com/robbymilo/rgallery/pkg/server"
	"github.com/robbymilo/rgallery/pkg/static"
)

func SetupRouter(c Conf, cache *cache.Cache, Commit, Tag string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RedirectSlashes)
	r.Use(middleware.Params(c))
	r.Use(middleware.Build(Commit, Tag))
	r.Use(middleware.Config(c))
	r.Use(middleware.Logger(c))

	// do not use metrics during test
	if !strings.HasSuffix(os.Args[0], ".test") {
		r.Use(chiprometheus.NewPatternMiddleware("rgallery"))
	}

	// serve files
	staticRoot, err := fs.Sub(static.StaticDir, ".")
	if err != nil {
		c.Logger.Error("error creating static filesystem", "error", err)
	}
	r.Handle("/static/*", Static(http.StripPrefix("/static/", http.FileServer(http.FS(staticRoot))), c))

	// handle fonts
	fontRoot, err := fs.Sub(fonts.FontDir, ".")
	if err != nil {
		c.Logger.Error("error creating static filesystem", "error", err)
	}
	r.Handle("/fonts/*", Fonts(http.StripPrefix("/fonts/", http.FileServer(http.FS(fontRoot))), c))

	r.Handle("/favicon.ico", http.FileServer(http.FS(staticRoot)))

	// thumbnails
	r.Route("/img", func(r chi.Router) {
		r.Use(middleware.Auth(c))
		r.Get("/{hash}/{size}", server.ServeThumbnail)
	})

	r.Route("/transcode", func(r chi.Router) {
		r.Use(middleware.Auth(c))
		r.Use(middleware.Logger(c))
		r.Get("/{hash}/{file}", server.ServeTranscode)
	})

	r.Route("/dist", func(r chi.Router) {
		r.Use(chiMiddleware.Compress(5))
		r.Use(middleware.Config(c))
		if c.Dev {
			// load files from dir for hot refresh
			fs := http.FileServer(http.Dir("./pkg/dist"))
			r.Handle("/*", http.StripPrefix("/dist/", fs))
		} else {
			// embed files
			distRoot, err := fs.Sub(dist.DistDir, ".")
			if err != nil {
				c.Logger.Error("error embeding dist dir", "error", err)
			}

			r.Handle("/*", Dist(http.StripPrefix("/dist/", http.FileServer(http.FS(distRoot))), c))
		}
	})

	r.Route("/", func(r chi.Router) {
		r.Use(chiMiddleware.Compress(5))
		r.Use(middleware.Config(c))
		r.Use(middleware.Cache(cache))
		r.Use(middleware.Auth(c))
		r.Use(middleware.Etag(c))

		// load originals from system
		r.Handle("/media-originals/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Strip the prefix
			path := strings.TrimPrefix(r.URL.Path, "/media-originals/")

			// Validate the path to prevent directory traversal
			if strings.Contains(path, "..") || strings.Contains(path, "//") {
				http.Error(w, "Invalid path", http.StatusBadRequest)
				return
			}

			// Serve the file
			http.StripPrefix("/media-originals/", http.FileServer(http.Dir(config.MediaPath(c)))).ServeHTTP(w, r)
		}))

		r.Get("/404", server.Send404)

		r.Get("/", server.ServeTimeline)
		r.Get("/onthisday", server.ServeOnThisDay)

		r.Get("/media/{hash}", server.ServeMedia)
		r.Get("/media/{hash}/in/{collection}", server.ServeMedia)         // for favorites
		r.Get("/media/{hash}/in/{collection}/{slug}*", server.ServeMedia) // for folders and tags

		r.Get("/folders", server.ServeFolders)
		r.Get("/folder*", server.ServeFolder)

		r.Get("/tags", server.ServeTags)
		r.Get("/tag/{slug}", server.ServeTag)

		r.Get("/favorites", server.ServeFavorites)
		r.Get("/map", server.ServeMap)
		r.Get("/gear", server.ServeGear)
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			server.ServeAdmin(w, r, c)
		})

		r.Get("/tiles/{z}/{x}/{y}.png", func(w http.ResponseWriter, r *http.Request) {
			server.ServeTiles(w, r, c)
		})

		r.NotFound(server.NotFound)
		r.MethodNotAllowed(server.NotAllowed)

		// scan only new/modified items
		r.Get("/scan", func(w http.ResponseWriter, r *http.Request) {
			server.Scan(w, r, "default", cache)
		})

		// remove and add all items
		// recreate thumbnails
		r.Get("/deepscan", func(w http.ResponseWriter, r *http.Request) {
			server.Scan(w, r, "deep", cache)
		})

		// remove and add all items
		// do not recreate thumbnails
		r.Get("/metadatascan", func(w http.ResponseWriter, r *http.Request) {
			server.Scan(w, r, "metadata", cache)
		})

		// check for missing thumbnails and generate missing ones
		// ignores the pregenerate-thumbs flag
		r.Get("/thumbscan", func(w http.ResponseWriter, r *http.Request) {
			server.ThumbScan(w, r)
		})

		r.Get("/poll", server.ServePoll)
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			server.ServeStatus(w, r, c)
		})

		if !c.DisableAuth {
			r.Get("/adduser", server.ServeAddUser)
			r.Get("/logout", server.ServeLogOut)
			r.Post("/admin/keys/create", server.CreateKey)
			r.Post("/admin/keys/delete", server.RemoveKey)

		}

	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok\n"))
		if err != nil {
			c.Logger.Error("error writing health check status", "error", err)
		}
	})

	if !c.DisableAuth {
		r.Get("/signin", server.ServeSignIn)
		r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
			err := server.SignIn(w, r, c)
			if err != nil {
				c.Logger.Error("error on signin route", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("503\n"))
				if err != nil {
					c.Logger.Error("error writing 500 status for signin route", "error", err)
				}
			}
		})
		r.Post("/signup", server.SignUp)

	}

	if c.ResizeService != "" {

		operation := func() error {
			c.Logger.Info("attempting to connect to resize service at " + c.ResizeService + "...")

			return server.CheckResizeServiceHealth(c)
		}

		_ = backoff.Retry(operation, backoff.NewExponentialBackOff())
		c.Logger.Info("connected to resize service at " + c.ResizeService)

	}

	if c.LocationService != "" {

		operation := func() error {
			c.Logger.Info("attempting to connect to location service at " + c.LocationService + "...")

			return server.CheckLocationServiceHealth(c)
		}

		_ = backoff.Retry(operation, backoff.NewExponentialBackOff())
		c.Logger.Info("connected to location service at " + c.LocationService)

	}

	return r

}

func SetupMetrics(c Conf) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Config(c))
	r.Use(middleware.Logger(c))

	// do not use metrics during test
	if !strings.HasSuffix(os.Args[0], ".test") {
		cs := metrics.MetricsCollector(c)
		prometheus.MustRegister(cs)

		r.Handle("/metrics", promhttp.Handler())
	}

	return r
}
