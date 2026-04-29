package main

import (
	"log"
	"net/http"
)

func main() {
	srv := setup()

	if err := srv.start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
