package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/robbymilo/rgallery/pkg/middleware"
	"github.com/robbymilo/rgallery/pkg/resize"
	"github.com/robbymilo/rgallery/pkg/types"
	cli "github.com/urfave/cli/v2"
)

type Conf = types.Conf

func main() {
	var tz string
	tz = os.Getenv("TZ")
	if tz == "" {
		tz = "UTC"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatal("error getting timezone:", err)
	}
	time.Local = loc

	port := "3001"

	var c = Conf{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	app := &cli.App{
		Name:  "rgallery-resize",
		Flags: []cli.Flag{},
		Action: func(cCtx *cli.Context) error {
			r := chi.NewRouter()
			r.Use(middleware.Logger(c))

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				err := resize.ResizeImageUpload(w, r, c)
				if err != nil {
					c.Logger.Error("error resizing upload", "error", err)
				}
			})

			r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("ok\n"))
				if err != nil {
					c.Logger.Error("error writing health check response", "error", err)
				}
			})

			c.Logger.Info("rgallery-resize listening on: " + port)
			err := http.ListenAndServe(":"+port, r)
			if err != nil {
				c.Logger.Error("error starting resizer", "error", err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
