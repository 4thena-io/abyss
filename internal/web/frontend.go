package web

import (
	"log"
	"net/http"
	"strings"

	"git.4thena.io/4thena/abys/frontend"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	fs, err := frontend.HttpFs()
	if err != nil {
		log.Fatalf("failed to load frontend assets: %v", err)
	}

	fileServer := http.FileServer(fs)

	spaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !strings.Contains(path, ".") || strings.HasSuffix(path, "/") {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	router.PathPrefix("/").Handler(spaHandler).Methods("GET")
}
