package main

import (
	"log"

	"github.com/rafaelbeecker/mwskit/internal/cmd/toolkit"
)

func main() {
	if err := toolkit.NewRootCmd().Execute(); err != nil {
		log.Fatalln(err)
	}
}
