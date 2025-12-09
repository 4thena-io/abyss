package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var app embed.FS

func HttpFs() (http.FileSystem, error) {
	dist, err := fs.Sub(app, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(dist), nil
}
