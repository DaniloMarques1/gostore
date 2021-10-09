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
	InvalidSyntax = "Message is not valid. Invalid syntax"
	// TODO better error report, it is not the key of the value it's the op, key, value
	InvalidMessageKey     = "The key given is not valid"
	OperationNotSupported = "Operation is not supported"
)

type StorageInterface interface {
	Store(key string, value string)
	Read(key string) string
	Delete(key string)
}

type Storage struct {
	db map[string]string
}

func NewStorage() *Storage {
	db := make(map[string]string)
	return &Storage{
		db: db,
	}
}

func (s *Storage) Store(key string, value string) {
	s.db[key] = value
}

func (s *Storage) Read(key string) string {
	return s.db[key]
}

func (s *Storage) Delete(key string) {
	delete(s.db, key)
}

type Message struct {
	op    string
	key   string
	value string
}

type Response struct {
	code int
	msg  string
}

type Operation interface {
	executeOperation(conn net.Conn, storage StorageInterface) error
}

type StoreOperation struct {
	key   string
	value string
}

type ReadOperation struct {
	key string
}

type DeleteOperation struct {
	key string
}

// we do use conn because we do not need send message to client
func (sop StoreOperation) executeOperation(_ net.Conn, storage StorageInterface) error {
	log.Printf("Executing store operation\n")
	storage.Store(sop.key, sop.value)
	return nil
}

// we do use conn because we do not need send message to client
func (dop DeleteOperation) executeOperation(_ net.Conn, storage StorageInterface) error {
	log.Printf("Executing delete operation\n")
	storage.Delete(dop.key)
	return nil
}

// we write to conn the read result
func (rop ReadOperation) executeOperation(conn net.Conn, storage StorageInterface) error {
	log.Printf("Executing read operation\n")
	value := storage.Read(rop.key)
	conn.Write([]byte(value))

	return nil
}

func main() {
	storage := NewStorage()
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

		go handleConnection(conn, storage)
	}
}

// handle the new connection
func handleConnection(conn net.Conn, storage *Storage) {
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
			conn.Write([]byte(err.Error() + "\n"))
			continue
		}
		op, err := getOperationFromMessage(message)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			continue
		}
		op.executeOperation(conn, storage)
	}
}

// parses the message to extract its key and value
// message will be something like op=store;key=swee;value=1234\n
func parseMessage(msg string) (*Message, error) {
	log.Printf("IN parseMessage %+v\n", msg)
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
	if opSplit[0] != "op" || len(opSplit[1]) == 0 {
		log.Printf("Invalid syntax on operation key\n")
		return nil, errors.New(InvalidMessageKey)
	}

	op := opSplit[1]

	keySplit := strings.Split(splited[1], "=")
	if len(keySplit) < 2 {
		log.Printf("Invalid syntax on key\n")
		return nil, errors.New(InvalidSyntax)
	}
	if keySplit[0] != "key" || len(keySplit[1]) == 0 {
		log.Printf("Invalid syntax on key key\n")
		return nil, errors.New(InvalidMessageKey)
	}
	key := keySplit[1]

	var value string
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
		if valueSplit[0] != "value" || len(valueSplit[1]) == 0 {
			log.Printf("Invalid syntax on key value.\n")
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
	log.Printf("IN getOperationFromMessage %+v\n", message)
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
