package main

import (
	//"encoding/json"
	"log"

	"github.com/danilomarques1/gostore/example/lib/golib"
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
	deleteOperation()
	keysOperation()
	listOperation()
}

func storeOperation() {
	resp, err := c.StoreOperation("isHidden", true)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", resp)
}

func readOperation() {
	response, err := c.ReadOperation("name")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", response)
}

func listOperation() {
	response, err := c.ListOperation()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", response)
}

func deleteOperation() {
	response, err := c.DeleteOperation("isHidden")
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
