package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	OP_STORE  = "store"
	OP_READ   = "read"
	OP_DELETE = "delete"
)

type Client struct {
	Host string
	Port string
	conn net.Conn
}

func NewClient(host, port string) *Client {
	return &Client{
		Host: host,
		Port: port,
	}
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%v:%v", c.Host, c.Port))
	if err != nil {
		log.Printf("Error reading message %v\n", err)
		return err
	}

	return nil
}

func (c *Client) Disconnect() error {
	return c.conn.Close()
}

func (c *Client) sendMessage(op Operation) (*Response, error) {
	bytes := op.parseOperation()
	c.conn.Write(bytes)
	responseString, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading message %v\n", err)
		return nil, err
	}
	response := parseResponse(responseString)

	return response, nil
}

type Operation interface {
	parseOperation() []byte
}

type storeOperation struct {
	key   string
	value interface{}
}

func (sop *storeOperation) parseOperation() []byte {
	s := fmt.Sprintf("op=store;key=%v;value=%v;\n", sop.key, sop.value)

	return []byte(s)
}

func (c *Client) StoreOperation(key string, value interface{}) (*Response, error) {
	op := storeOperation{key: key, value: value}
	resp, err := c.sendMessage(&op)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ReadOperation(key string) (*Response, error) {
	op := readOperation{key: key}
	resp, err := c.sendMessage(&op)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type readOperation struct {
	key string
}

func (rop *readOperation) parseOperation() []byte {
	s := fmt.Sprintf("op=read;key=%v;\n", rop.key)

	return []byte(s)
}

type deleteOperation struct {
	key string
}

func (c *Client) DeleteOperation(key string) (*Response, error) {
	op := deleteOperation{key: key}
	resp, err := c.sendMessage(&op)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (dop *deleteOperation) parseOperation() []byte {
	s := fmt.Sprintf("op=delete;key=%v;\n", dop.key)

	return []byte(s)
}

type listOperation struct {
}

func (lop *listOperation) parseOperation() []byte {
	return []byte("op=list;\n")
}

func (c *Client) ListOperation() (*Response, error) {
	op := listOperation{}
	resp, err := c.sendMessage(&op)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type Response struct {
	Code    int
	Message interface{}
}

func NewResponse(code int, message interface{}) *Response {
	return &Response{Code: code, Message: message}
}

func parseResponse(responseStr string) *Response {
	splited := strings.Split(responseStr, ";")
	codeSplit := strings.Split(splited[0], "=")
	code, _ := strconv.Atoi(codeSplit[1]) // TODO
	msgSplit := strings.Split(splited[1], "=")
	msg := msgSplit[1]

	return NewResponse(code, msg)
}
