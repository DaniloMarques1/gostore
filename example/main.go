package main

import (
	"log"

	"github.com/danilomarques1/gostore/example/client"
)

const (
	HOST = "127.0.0.1"
	PORT = "5000"
)

func main() {
	c := client.NewClient(HOST, PORT)
	if err := c.Connect(); err != nil {
		log.Fatalf("ERR connecting to server %v\n", err)
	}

	msg := client.NewMessage(client.OP_STORE, "name", "Danilo")
	response, err := c.SendMessage(msg)
	if err != nil {
		log.Fatalf("Error sending message %v\n", err)
	}
	log.Printf("Response = %v\n", string(response))

	msg = client.NewMessage(client.OP_READ, "name", "")
	response, err = c.SendMessage(msg)
	if err != nil {
		log.Fatalf("ERROR sending message %v\n", err)
	}

	log.Printf("Response = %v\n", string(response))
}
