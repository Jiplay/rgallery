package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/robbymilo/rgallery/pkg/geo"
	"github.com/robbymilo/rgallery/pkg/middleware"
	"github.com/robbymilo/rgallery/pkg/types"
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

	var c = Conf{
		LocationDataset: "Provinces10",
		Logger:          slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	port := "3002"
	r := chi.NewRouter()
	r.Use(middleware.Logger(c))

	h, err := geo.NewGeoHandler(c)
	if err != nil {
		c.Logger.Error("error getting new handlers", "error", err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		lon := r.URL.Query().Get("lon")
		lat := r.URL.Query().Get("lat")

		longitude, err := strconv.ParseFloat(lon, 64)
		if err != nil {
			c.Logger.Error("error parsing longitude", "error", err)
		}
		latitude, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			c.Logger.Error("error parsing latitude", "error", err)
		}

		loc, err := geo.GetLocation(h, longitude, latitude, c)
		if err != nil {
			c.Logger.Error("error getting location", "error", err)
		}

		json, err := json.Marshal(loc)
		if err != nil {
			c.Logger.Error("error marshalling json", "error", err)
		}

		_, err = w.Write(json)
		if err != nil {
			c.Logger.Error("error writing json", "error", err)
		}

	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok\n"))
		if err != nil {
			c.Logger.Error("error writing health check response", "error", err)
		}

	})

	c.Logger.Info("rgallery-geo listening on: " + port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		c.Logger.Error("error starting geo", "error", err)
	}

}
