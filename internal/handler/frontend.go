package handler

import (
	"log"
	"net/http"
	"strings"

	"git.d4ramirez.com/project-abyss/abys-api/frontend"
	"github.com/gorilla/mux"
)

type FrontendHandler struct {
	staticFileHandler http.Handler
}

func NewFrontendHandler() *FrontendHandler {
	uiFs, err := frontend.HttpFs()
	if err != nil {
		log.Fatal("Failed to load frontend assets: %w", err)
	}

	fileServer := http.FileServer(uiFs)

	spaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if !strings.Contains(p, ".") || strings.HasSuffix(p, "/") {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	return &FrontendHandler{
		staticFileHandler: spaHandler,
	}
}

func (h *FrontendHandler) RegisterFrontendRoutes(router *mux.Router) {
	router.PathPrefix("/").Handler(h.staticFileHandler).Methods("GET")
}
