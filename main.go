package main

import (
	"flag"
	"log"

	"github.com/danilomarques1/gostore/server"
)

var port string

func main() {
	parseFlag()
	s, err := server.NewServer(port)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}

func parseFlag() {
	flag.StringVar(&port, "port", "5000", "The server port")
	flag.Parse()
}
