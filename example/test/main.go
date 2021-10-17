package main

import (
	"log"

	"github.com/danilomarques1/gostore/example/client"
)

var c *client.Client

func main() {
	c = client.NewClient("localhost", "5000")
	if err := c.Connect(); err != nil {
		log.Fatalf("ERR %v\n", err)
	}
	defer c.Disconnect()

	storeOperation()
	readOperation()
	listOperation()
	keysOperation()
	deleteOperation()
	listOperation()
}

func storeOperation() {
	resp, err := c.StoreOperation("name", "Danilo")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", resp)
}

func listOperation() {
	response, err := c.ListOperation()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", response)
}

func readOperation() {
	response, err := c.ReadOperation("name")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", response)
}

func deleteOperation() {
	response, err := c.DeleteOperation("name")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", response)
}

func keysOperation() {
	response, err := c.KeysOperation()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", response)
}
