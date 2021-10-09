package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strings"
)

const (
	// TODO DEFAULT PORT
	ADDRESS = ":5000"
)

// operations supported
const (
	OP_STORE  = "store"
	OP_DELETE = "delete"
	OP_READ   = "read"
)

// error messages
const (
	InvalidSyntax         = "Message is not valid. Invalid syntax"
	OperationNotSupported = "Operation is not supported"
)

type Store struct {
	db map[string]interface{}
}

type Message struct {
	op    string
	key   string
	value interface{}
}

type Response struct {
	code int
	msg  string
}

type Operation interface {
	executeOperation() error
}

type StoreOperation struct {
	key   string
	value interface{}
}

type ReadOperation struct {
	key string
}

type DeleteOperation struct {
	key string
}

func (sop StoreOperation) executeOperation() error {
	return nil
}

func (dop DeleteOperation) executeOperation() error {
	return nil
}

func (rop ReadOperation) executeOperation() error {
	return nil
}

func main() {
	socket, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		log.Fatalf("ERR: Error opening tcp connection %v\n", err)
	}
	log.Printf("Starting server...\n")
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("ERR: Error accepting connection %v\n", err)
			continue
		}
		log.Printf("New connection arrived\n")

		go handleConnection(conn)
	}
}

// handle the new connection
func handleConnection(conn net.Conn) {
	log.Printf("Handling connection\n")
	for {
		//defer conn.Close()

		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("ERR: Error reading message from connection %v\n", err)
			return
		}

		message, err := parseMessage(msg)
		if err != nil {
			//TODO
		}
		op, err := getOperationFromMessage(message)
		if err != nil {
			// TODO handl err
		}
		op.executeOperation()
	}
}

// parses the message to extract its key and value
// message will be something like op=store;key=swee;value=1234\n
func parseMessage(msg string) (*Message, error) {
	msg = strings.Replace(msg, "\n", "", -1)
	splited := strings.Split(msg, ";")
		
	if len(splited) < 2 {
		log.Printf("Invalid syntax. Wrong number of ;\n")
		return nil, errors.New(InvalidSyntax)
	}

	opSplit := strings.Split(splited[0], "=")
	if len(opSplit) < 2 {
		log.Printf("Invalid syntax on operation\n")
		return nil, errors.New(InvalidSyntax)
	}
	op := opSplit[1]

	keySplit := strings.Split(splited[1], "=")
	if len(keySplit) < 2 {
		log.Printf("Invalid syntax on key\n")
		return nil, errors.New(InvalidSyntax)
	}
	key := keySplit[1]

	var value interface{}
	if op == OP_STORE {
		if len(splited) < 3 {
			log.Printf("Invalid syntax on value. You did not provide a value key-value pair\n")
			return nil, errors.New(InvalidSyntax)
		}
		valueSplit := strings.Split(splited[2], "=")
		log.Printf("Value split = %v\n", valueSplit)
		if len(valueSplit) < 2 {
			log.Printf("Invalid syntax on value.\n")
			return nil, errors.New(InvalidSyntax)
		}
		value = valueSplit[1]
	}

	return &Message{
		op:    op,
		key:   key,
		value: value,
	}, nil
}

// return the Operation that will be executed
// based on the received message
func getOperationFromMessage(message *Message) (Operation, error) {
	var operation Operation
	if message.op == OP_STORE {
		operation = StoreOperation{
			key:   message.key,
			value: message.value,
		}
	} else if message.op == OP_DELETE {
		operation = DeleteOperation{
			key: message.key,
		}
	} else if message.op == OP_READ {
		operation = ReadOperation{
			key: message.key,
		}
	} else {
		return nil, errors.New(OperationNotSupported) 
	}

	return operation, nil
}
