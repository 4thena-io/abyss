package main

import (
	"log"

	"github.com/4thena-io/abyss/internal/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		log.Fatal(err)
	}
}
