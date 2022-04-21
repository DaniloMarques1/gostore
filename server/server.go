package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	listener net.Listener
	port     string
	storage  StorageInterface
}

type Message struct {
	op    string
	key   string
	value interface{}
}

// server error messages
const (
	InvalidSyntax = "Message is not valid. Invalid syntax"
	// TODO better error report, it is not the key of the value it's the op, key, value
	InvalidMessageKey           = "The given key is not valid"
	OperationNotSupported       = "Operation is not supported"
	KeyNotFound                 = "The given key was not found"
	DuplicationOfKey            = "The given key is already registered"
	OperationOnTypeNotSupported = "The value type does not support this operation"
)

func NewServer(port string) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}
	storage := NewStorage()
	return &Server{listener: listener, port: port, storage: storage}, nil
}

func (s *Server) Start() {
	log.Printf("Starting server on port %v\n", s.port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("New connection\n")
		go s.handleConnection(conn) // TODO add clover db
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Printf("Handling connection: %v\n", conn.RemoteAddr().String())
	defer conn.Close()
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("ERR: Error reading message from connection %v\n", err)
			return
		}

		message, err := s.parseMessage(msg)
		if err != nil {
			conn.Write([]byte(s.parseResponse(1, err.Error())))
			continue
		}
		op, err := s.getOperationFromMessage(message)
		if err != nil {
			conn.Write([]byte(s.parseResponse(1, err.Error())))
			continue
		}

		responseValue, err := op.ExecuteOperation(s.storage)
		if err != nil {
			conn.Write([]byte(s.parseResponse(1, err.Error())))
			continue
		}

		conn.Write([]byte(s.parseResponse(0, responseValue)))
	}
}

// parseMessage convert the received string into a Message
func (s *Server) parseMessage(msg string) (*Message, error) {
	log.Printf("parseMessage %+v\n", msg)
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
	if op == OP_STORE || op == OP_REPLACE {
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

// getOperationFromMessage
// return the Operation that will be executed
// based on the received message
func (s *Server) getOperationFromMessage(m *Message) (Operation, error) {
	log.Printf("IN getOperationFromMessage %+v\n", m)
	var operation Operation
	switch m.op {
	case OP_STORE:
		operation = StoreOperation{key: m.key, value: m.value}
	case OP_DELETE:
		operation = DeleteOperation{key: m.key}
	case OP_READ:
		operation = ReadOperation{key: m.key}
	case OP_LIST:
		operation = ListOperation{}
	case OP_KEYS:
		operation = KeysOperation{}
	case OP_REPLACE:
		operation = ReplaceOperation{key: m.key, value: m.value}
	default:
		return nil, errors.New(OperationNotSupported)
	}

	return operation, nil
}

// parseResponse format the status and the message into a response to the client 
func (s *Server) parseResponse(code int, message interface{}) string {
	return fmt.Sprintf("code=%v;message=%v\n", code, message)
}
