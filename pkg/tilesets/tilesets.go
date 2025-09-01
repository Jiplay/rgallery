package tilesets

import (
	"embed"
)

//go:embed *
var TileSets embed.FS

func ExposeEmbeddedFile() ([]byte, error) {

	return TileSets.ReadFile("world.mbtiles")
}
