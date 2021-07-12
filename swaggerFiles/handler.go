package swaggerFiles

import (
	"embed"
	"io/fs"
	"net/http"
)

// It will add all the files in dist
//go:embed dist/*
var static embed.FS

var Handler http.Handler

func init() {
	contentFS, _ := fs.Sub(static, "dist")
	Handler = http.FileServer(http.FS(contentFS))
}
