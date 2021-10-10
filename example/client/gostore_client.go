package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
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

func (c *Client) SendMessage(m *Message) ([]byte, error) {
	bytes := m.parseMessage()
	c.conn.Write(bytes)
	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading message %v\n", err)
		return []byte(""), err
	}

	return []byte(response), nil
}

type Message struct {
	op    string
	key   string
	value interface{}
}

func NewMessage(op, key, value string) *Message {
	return &Message{
		op:    op,
		key:   key,
		value: value,
	}
}

func (m *Message) parseMessage() []byte {
	var b []byte
	if m.op == OP_STORE {
		b = []byte(fmt.Sprintf("op=%v;key=%v;value=%v\n",
			m.op, m.key, m.value))
	} else {
		b = []byte(fmt.Sprintf("op=%v;key=%v\n",
			m.op, m.key))
	}

	return b
}
