package config

import (
	"html/template"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/robbymilo/rgallery/pkg/types"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type Meta = types.Meta
type Media = types.Media
type Conf = types.Conf

func CachePath(c Conf) string {
	cache_path, err := filepath.Abs(c.Cache)
	if err != nil {
		c.Logger.Error("error parsing cache dir", "error", err)
	}

	return cache_path
}

func MediaPath(c Conf) string {
	media_path, err := filepath.Abs(c.Media)
	if err != nil {
		c.Logger.Error("error parsing media dir", "error", err)
	}

	return media_path
}

// GetConf returns a Conf struct from the config file, cli flags, and env vars.
func GetConf(cCtx cli.Context, Commit, Tag string) Conf {
	var c Conf

	data, err := os.ReadFile(cCtx.String("config"))
	if err != nil {
		c.Logger.Info("Optional config file not found", "error", err)
		return c
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		c.Logger.Error("Error parsing YAML file", "error", err)
		return c
	}

	c = Conf{
		Dev:              cCtx.Bool("dev"),
		DisableAuth:      cCtx.Bool("disable-auth"),
		Media:            cCtx.String("media"),
		Cache:            cCtx.String("cache"),
		Data:             cCtx.String("data"),
		Quality:          cCtx.Int("quality"),
		PreGenerateThumb: cCtx.Bool("pregenerate-thumbs"),
		ResizeService:    cCtx.String("resize_service"),
		LocationService:  cCtx.String("location-service"),
		LocationDataset:  cCtx.String("location-dataset"),
		Logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
		TileServer:       cCtx.String("tile-server"),
		SessionLength:    cCtx.Int("session-length"),
		IncludeOriginals: cCtx.Bool("include-originals"),
		Meta: Meta{
			Commit:     Commit,
			Tag:        Tag,
			CustomHTML: template.HTML(c.CustomHTML),
		},
		OnThisDay: cCtx.Bool("on-this-day"),
	}

	return c
}
