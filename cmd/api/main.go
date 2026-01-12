package main

import (
	"log"
	"net/http"

	"git.4thena.io/4thena/abys/internal/server"
)

func main() {
	srv := server.New(server.Config{
		Host: "0.0.0.0",
		Port: "8000",
	})

	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
