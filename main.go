package main

import (
	"log"

	"github.com/4thena-io/abyss/internal/cli"
)

func main() {
	if err := cli.Root().Execute(); err != nil {
		log.Fatal(err)
	}
}
