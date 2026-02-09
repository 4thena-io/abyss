package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/4thena-io/abyss/frontend"
)

type FrontendHandler struct {
	fileServer http.Handler
}

func NewFrontendHandler() *FrontendHandler {
	uiFs, err := frontend.HttpFs()
	if err != nil {
		log.Fatalf("failed to load frontend assets: %v", err)
	}

	return &FrontendHandler{
		fileServer: http.FileServer(uiFs),
	}
}

func (h *FrontendHandler) Serve(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.Contains(path, ".") || strings.HasSuffix(path, "/") {
		r.URL.Path = "/"
	}
	h.fileServer.ServeHTTP(w, r)
}
