package main

import (
	"log"
	"net/http"

	"git.4thena.io/4thena/abys/internal/handler"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	handler.NewAppHandler().RegisterAppRoutes(apiRouter)
	handler.NewProjectHandler().RegisterProjectRoutes(apiRouter)

	handler.NewFrontendHandler().RegisterFrontendRoutes(router)

	var host = "0.0.0.0:8000"

	log.Printf("Running on http://%s\n", host)

	if err := http.ListenAndServe(host, router); err != nil {
		log.Fatal("Failed to start http server")
	}
}
