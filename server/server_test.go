package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"testing"

	c "github.com/ostafen/clover"
	"github.com/stretchr/testify/require"
)

const DIR = "test-db"

func createDir() error {
	return os.Mkdir(DIR, os.ModePerm)
}

func removeDir() error {
	return os.RemoveAll(DIR)
}

func getStorage() *Storage {
	db, _ := c.Open(DIR)
	collectionName := "gostore-test"
	db.CreateCollection(collectionName)
	storage := &Storage{db: db, collectionName: collectionName }
	return storage
}

func TestStore(t *testing.T) {
	require.NoError(t, createDir())
	defer removeDir()
	storage := getStorage()
	s := &Server{port: "8080", storage: storage}

	serverConn, clientConn := net.Pipe()
	go func(s *Server) {
		s.handleConnection(serverConn)
		serverConn.Close()
	}(s)
	_, err := clientConn.Write([]byte("op=store;key=name;value=Danilo\n"));
	require.NoError(t, err)

	msg, err := bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, "code=0;message=Value stored successfully\n", msg)

	clientConn.Close()
}

func TestRead(t *testing.T) {
	require.NoError(t, createDir())
	defer removeDir()
	storage := getStorage()
	s := &Server{port: "8080", storage: storage}

	serverConn, clientConn := net.Pipe()
	go func(s *Server) {
		s.handleConnection(serverConn)
		serverConn.Close()
	}(s)
	_, err := clientConn.Write([]byte("op=store;key=name;value=Danilo\n"));
	require.NoError(t, err)
	msg, err := bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, "code=0;message=Value stored successfully\n", msg)

	_, err = clientConn.Write([]byte("op=read;key=name\n"))
	require.NoError(t, err)

	msg, err = bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, "code=0;message=Danilo\n", msg)

	clientConn.Close()
}

func TestDelete(t *testing.T) {
	require.NoError(t, createDir())
	defer removeDir()
	storage := getStorage()
	s := &Server{port: "8080", storage: storage}

	serverConn, clientConn := net.Pipe()
	go func(s *Server) {
		s.handleConnection(serverConn)
		serverConn.Close()
	}(s)
	_, err := clientConn.Write([]byte("op=store;key=name;value=Danilo\n"));
	require.NoError(t, err)
	msg, err := bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("code=0;message=%v\n", StoredSuccessFully), msg)

	_, err = clientConn.Write([]byte("op=delete;key=name;\n"))
	require.NoError(t, err)
	msg, err = bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("code=0;message=%v\n", DeletedSuccessFully), msg)


	_, err = clientConn.Write([]byte("op=read;key=name\n"))
	require.NoError(t, err)
	msg, err = bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("code=1;message=%v\n", KeyNotFound), msg) 

	clientConn.Close()
}

func TestKeys(t *testing.T) {
	require.NoError(t, createDir())
	defer removeDir()
	storage := getStorage()
	s := &Server{port: "8080", storage: storage}

	serverConn, clientConn := net.Pipe()
	go func(s *Server) {
		s.handleConnection(serverConn)
		serverConn.Close()
	}(s)
	_, err := clientConn.Write([]byte("op=store;key=name;value=Danilo\n"))
	require.NoError(t, err)

	msg, err := bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("code=0;message=%v\n", StoredSuccessFully), msg)

	_, err = clientConn.Write([]byte("op=keys\n"))
	require.NoError(t, err)

	msg, err = bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, "code=0;message=[name]\n", msg)

	clientConn.Close()
}

func TestReplace(t *testing.T) {
	require.NoError(t, createDir())
	defer removeDir()
	storage := getStorage()
	s := &Server{port: "8080", storage: storage}

	serverConn, clientConn := net.Pipe()
	go func(s *Server) {
		s.handleConnection(serverConn)
		serverConn.Close()
	}(s)
	_, err := clientConn.Write([]byte("op=store;key=name;value=Danilo\n"))
	require.NoError(t, err)

	msg, err := bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("code=0;message=%v\n", StoredSuccessFully), msg)

	_, err = clientConn.Write([]byte("op=read;key=name\n"))
	require.NoError(t, err)

	msg, err = bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, "code=0;message=Danilo\n", msg)

	_, err = clientConn.Write([]byte("op=replace;key=name;value=Fitz\n"))
	require.NoError(t, err)

	msg, err = bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("code=0;message=%v\n", ReplacedSuccessFully), msg)

	_, err = clientConn.Write([]byte("op=read;key=name\n"))
	require.NoError(t, err)

	msg, err = bufio.NewReader(clientConn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, "code=0;message=Fitz\n", msg)

	clientConn.Close()
}

