package templates

import "embed"

//go:embed */*.html
var TemplatesDir embed.FS
