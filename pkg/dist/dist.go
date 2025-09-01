package dist

import "embed"

//go:embed */*
var DistDir embed.FS
