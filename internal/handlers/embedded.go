package handlers

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

// GetStaticFileSystem returns a http.FileSystem for serving static files
func GetStaticFileSystem() (http.FileSystem, error) {
	// Strip the "static" prefix from the embedded files
	stripped, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, err
	}
	return http.FS(stripped), nil
}
