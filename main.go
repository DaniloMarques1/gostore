package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	// TODO DEFAULT PORT
	ADDRESS = ":5000"
)

// error messages
const (
	InvalidSyntax = "Message is not valid. Invalid syntax"
	// TODO better error report, it is not the key of the value it's the op, key, value
	InvalidMessageKey           = "The given key is not valid"
	OperationNotSupported       = "Operation is not supported"
	KeyNotFound                 = "The given key was not found"
	DuplicationOfKey            = "The given key is already registered"
	OperationOnTypeNotSupported = "The value type does not support this operation"
)

// Response Messages
const (
	StoredSuccessFully  = "Value stored successfully"
	DeletedSuccessFully = "Value removed successfully"
)

type Message struct {
	op    string
	key   string
	value interface{}
}

type Response struct {
	code    int
	message interface{}
}

func NewResponse(code int, message interface{}) *Response {
	return &Response{
		code:    code,
		message: message,
	}
}

func (r *Response) parseResponse() string {
	return fmt.Sprintf("code=%v;message=%v\n", r.code, r.message)
}

func main() {
	storage := NewStorage()
	SyncRead(&storage.db)

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
	defer conn.Close()
	for {
		log.Printf("DB = %+v\n", storage.db)
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("ERR: Error reading message from connection %v\n", err)
			return
		}

		message, err := parseMessage(msg)
		if err != nil {
			response := NewResponse(1, err.Error())
			conn.Write([]byte(response.parseResponse()))
			continue
		}
		op, err := getOperationFromMessage(message)
		if err != nil {
			response := NewResponse(1, err.Error())
			conn.Write([]byte(response.parseResponse()))
			continue
		}

		responseValue, err := op.ExecuteOperation(storage)
		if err != nil {
			response := NewResponse(1, err.Error())
			conn.Write([]byte(response.parseResponse()))
			continue
		}

		if op.GetOpType() == OP_STORE || op.GetOpType() == OP_DELETE {
			go SyncWrite(&storage.db)
		}

		response := NewResponse(0, responseValue)
		conn.Write([]byte(response.parseResponse()))
	}
}

// parses the message to extract its key and value
// message will be something like op=store;key=swee;value=1234\n
func parseMessage(msg string) (*Message, error) {
	log.Printf("IN parseMessage %+v\n", msg)
	msg = strings.Replace(msg, "\n", "", -1)
	msg = strings.Replace(msg, "\r", "", -1)
	splited := strings.Split(msg, ";")

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

	if (op != OP_LIST && op != OP_KEYS) && len(splited) < 2 {
		log.Printf("Invalid syntax. Wrong number of ;\n")
		return nil, errors.New(InvalidSyntax)
	}

	var key string
	if op != OP_LIST && op != OP_KEYS {
		keySplit := strings.Split(splited[1], "=")
		if len(keySplit) < 2 {
			log.Printf("Invalid syntax on key\n")
			return nil, errors.New(InvalidSyntax)
		}
		if keySplit[0] != "key" || len(keySplit[1]) == 0 {
			log.Printf("Invalid syntax on key key\n")
			return nil, errors.New(InvalidMessageKey)
		}

		key = keySplit[1]
	}

	var value string
	if op == OP_STORE {
		if len(splited) < 3 {
			log.Printf("Invalid syntax on value. You did not provide a value for the key-value pair\n")
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
func getOperationFromMessage(m *Message) (Operation, error) {
	log.Printf("IN getOperationFromMessage %+v\n", m)
	var operation Operation
	if m.op == OP_STORE {
		operation = StoreOperation{
			key:   m.key,
			value: m.value,
		}
	} else if m.op == OP_DELETE {
		operation = DeleteOperation{
			key: m.key,
		}
	} else if m.op == OP_READ {
		operation = ReadOperation{
			key: m.key,
		}
	} else if m.op == OP_LIST {
		operation = ListOperation{}
	} else if m.op == OP_KEYS {
		operation = KeysOperation{}
	} else {
		return nil, errors.New(OperationNotSupported)
	}

	return operation, nil
}

func SyncRead(db *map[string]interface{}) {
	file, err := os.OpenFile("db", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err) // TODO better report
	}
	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err) // TODO
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, db)
		if err != nil {
			log.Fatal(err) // TODO
		}
	}

	file.Close()
}

func SyncWrite(db *map[string]interface{}) {
	file, err := os.OpenFile("db", os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err) // TODO better report
	}
	b, err := json.Marshal(db)
	if err != nil {
		log.Fatal(err) // TODO better report
	}

	_, err = file.Write(b)
	if err != nil {
		log.Fatal(err) // TODO
	}

	file.Close()
}
