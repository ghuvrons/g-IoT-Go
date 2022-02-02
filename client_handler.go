package giotgo

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type ClientHandler struct {
	connection *net.Conn
	bufReader  *bufio.Reader
	buffer     []byte
	timeout    int
}

func NewClientHandler(conn *net.Conn, timeout int) *ClientHandler {
	client := &ClientHandler{}
	client.connection = conn
	client.bufReader = bufio.NewReader(*client.connection)
	client.buffer = make([]byte, 1024)
	client.timeout = timeout

	go func() {
		client.handle()
	}()
	return client
}

func (client *ClientHandler) handle() {
	fmt.Println("client handled")
	for true {
		if client.timeout != 0 {
			(*client.connection).SetReadDeadline(time.Now().Add(time.Duration(client.timeout) * time.Second))
		}
		b, err := client.bufReader.Read(client.buffer)
		if err != nil {
			fmt.Println((err))
			if err == io.EOF {
				break
			}
			continue
		}
		if b != 0 {
			fmt.Println(b, client.buffer)
		}
	}
}
