package giotgo

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

type ClientHandler struct {
	connection *net.Conn
	state      ClientState
	tmpBuffer  []byte
	buffer     *bytes.Buffer
	bufferPtr  int
	timeout    int
	tmpPacket  *packet
}

func NewClientHandler(conn *net.Conn, timeout int) *ClientHandler {
	client := &ClientHandler{}
	client.connection = conn
	client.tmpBuffer = make([]byte, 128)
	client.buffer = &bytes.Buffer{}
	client.timeout = timeout
	client.tmpPacket = NewPacket()

	go func() {
		client.state = CLIENT_STATE_CONNECTING
		client.handle()
	}()
	return client
}

func (client *ClientHandler) handle() {
	for true {
		if client.timeout != 0 {
			(*client.connection).SetReadDeadline(time.Now().Add(time.Duration(client.timeout) * time.Second))
		}

		len, err := (*client.connection).Read(client.tmpBuffer)

		if err != nil {
			fmt.Println((err))
			if err == io.EOF {
				break
			}
			continue
		}

		if len > 0 {
			client.buffer.Write(client.tmpBuffer[:len])
		}

		for client.buffer.Len() > 0 {
			if client.tmpPacket.packetType == 0 {
				_, err := client.tmpPacket.packetDecode(client.buffer)
				if err != nil {
					break
				}
			} else {
				willWriteLen := int(client.tmpPacket.length)
				if willWriteLen > client.buffer.Len() {
					willWriteLen = client.buffer.Len()
				}
				client.tmpPacket.payload.Write(client.buffer.Bytes()[:willWriteLen])
				client.buffer.Next(willWriteLen)
			}
		}
	}
}

func (client *ClientHandler) handlePacket() {

}
