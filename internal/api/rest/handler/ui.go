package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/4thena-io/abyss/internal/web"
)

type UIHandler struct {
	fileServer http.Handler
}

func NewUIHandler() *UIHandler {
	uiFs, err := web.HttpFs()
	if err != nil {
		log.Fatalf("failed to load UI assets: %v", err)
	}
	if uiFs == nil {
		return &UIHandler{fileServer: http.NotFoundHandler()}
	}

	return &UIHandler{
		fileServer: http.FileServer(uiFs),
	}
}

func (h *UIHandler) Serve(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.Contains(path, ".") || strings.HasSuffix(path, "/") {
		r.URL.Path = "/"
	}
	h.fileServer.ServeHTTP(w, r)
}
