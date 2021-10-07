package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

const (
	// TODO DEFAULT PORT
	ADDRESS = ":5000"
)

type Store struct {
	db map[string]interface{}
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

func handleConnection(conn net.Conn) {
	log.Printf("Handling connection\n")
	//defer conn.Close()

	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("ERR: Error reading message from connection %v\n", err)
		return
	}

	handleMessage(msg)
}

// parses the message to extract its key and value
// message will be something like key=swee;value=1234\n
func handleMessage(msg string) {
	msgWithoutNewLine := strings.Replace(msg, "\n", "", -1)

	splited := strings.Split(msgWithoutNewLine, ";")
	keySplit := strings.Split(splited[0], "=")
	key := keySplit[1]
	valueSplit := strings.Split(splited[1], "=")
	value := valueSplit[1]

	log.Printf("\nKEY = %v\nVALUE = %v\n", key, value)
}
