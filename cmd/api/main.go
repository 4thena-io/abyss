package main

import (
	"log"
	"net/http"

	"git.d4ramirez.com/project-abyss/abys-api/internal/app"
	"git.d4ramirez.com/project-abyss/abys-api/internal/project"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	project.NewHandler().RegisterRoutes(router)
	app.NewHandler().RegisterRoutes(router)

	var host = "0.0.0.0:8000"

	log.Printf("Running on http://%s\n", host)

	if err := http.ListenAndServe(host, router); err != nil {
		log.Fatal("Failed to start http server")
	}
}
