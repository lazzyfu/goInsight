package web

import (
	"embed"
)

// StaticFS embeds all static frontend files
//
//go:embed dist
var StaticFS embed.FS
