package main

import (
	"log"

	"github.com/danilomarques1/gostore/server"
)

func main() {
	s, err := server.NewServer("5000") // TODO add flag
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}
