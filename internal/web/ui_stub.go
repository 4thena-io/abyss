//go:build !embed_ui

package web

import "net/http"

// HttpFs is a stub for building Abyss without the embedded UI.
func HttpFs() (http.FileSystem, error) {
	return nil, nil
}
