package router

import (
	"net/http"
	"path/filepath"
)

// serves the HTML files with path and filename
func serveFile(basePath, file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(basePath, file))
	}
}
