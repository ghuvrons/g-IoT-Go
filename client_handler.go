package giotgo

import (
	"bytes"
	"io"
	"net"
	"time"

	"github.com/ghuvrons/g-IoT-Go/giot_packet"
)

type ClientHandler struct {
	connection   *net.Conn
	server       *Server
	state        ClientState
	tmpBuffer    []byte
	buffer       *bytes.Buffer
	bufferPtr    int
	timeout      int
	tmpPacket    *giot_packet.Packet
	queueCommand chan [2]interface{}

	info struct {
		name string
	}
}

func NewClientHandler(conn *net.Conn, server *Server, timeout int) *ClientHandler {
	client := &ClientHandler{}
	client.connection = conn
	client.tmpBuffer = make([]byte, 128)
	client.buffer = &bytes.Buffer{}
	client.timeout = timeout
	client.tmpPacket = giot_packet.NewPacket()
	client.server = server
	client.queueCommand = make(chan [2]interface{}, 5)

	go func() {
		client.state = CLIENT_STATE_CONNECTING
		client.handle()
	}()
	return client
}

func (client *ClientHandler) handle() {
	go client.handleCommand()

	for client.state != CLIENT_STATE_CLOSE {
		if client.timeout != 0 {
			(*client.connection).SetReadDeadline(time.Now().Add(time.Duration(client.timeout) * time.Second))
		}

		len, err := (*client.connection).Read(client.tmpBuffer)

		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		if len > 0 {
			client.buffer.Write(client.tmpBuffer[:len])
		}

		for client.buffer.Len() > 0 {
			if client.tmpPacket.PacketType == 0 {
				if err := client.tmpPacket.Decode(client.buffer); err != nil {
					break
				}
			} else {
				willWriteLen := int(client.tmpPacket.Length)
				if willWriteLen > client.buffer.Len() {
					willWriteLen = client.buffer.Len()
				}
				client.tmpPacket.Payload.Write(client.buffer.Bytes()[:willWriteLen])
				client.buffer.Next(willWriteLen)
			}

			if client.tmpPacket.IsValid() {
				client.handlePacket(client.tmpPacket)
				client.tmpPacket.Reset()
			}
		}
	}
}

func (client *ClientHandler) handlePacket(pck *giot_packet.Packet) {
	switch pck.PacketType {
	case giot_packet.PACKET_TYPE_CONNECT:
		pckConn := giot_packet.PacketConnectDecode(pck)

		// validating
		if err := pckConn.Validate(client.server.authenticator); err != nil {
			client.close()
		}

		// on success
		client.state = CLIENT_STATE_CONNECT

		bufConnack := &bytes.Buffer{}
		pckConnAck := giot_packet.NewPacketConnack(giot_packet.RESP_OK)
		pckConnAck.Encode(bufConnack)

		(*client.connection).Write(bufConnack.Bytes())

	case giot_packet.PACKET_TYPE_COMMAND:
		pckCmd := giot_packet.PacketCommandDecode(pck)

		var respStatus giot_packet.RespStatus = giot_packet.RESP_OK
		var respBuffer *bytes.Buffer

		handler, isOK := client.server.commandHandlers[pckCmd.Command]
		if !isOK && handler == nil {
			return
		} else {
			respBuffer = handler(client, pckCmd.Payload)
		}

		bufConnack := &bytes.Buffer{}
		packResp := giot_packet.NewPacketResponse(respStatus)
		packResp.AckId = pckCmd.AckId
		packResp.Payload = respBuffer
		packResp.Encode(bufConnack)

		(*client.connection).Write(bufConnack.Bytes())
	}
}

func (client *ClientHandler) handleCommand() {
	for client.state != CLIENT_STATE_CLOSE {
		select {
		case ch := <-client.queueCommand:
			exec, isOK := client.server.commandExecutors[ch[0].(giot_packet.Command)]
			if !isOK && exec == nil {
				continue
			} else {
				exec(client, ch[1])
			}

		case <-time.After(60 * time.Second):
			continue
		}

	}
}

func (client *ClientHandler) Execute(cmd giot_packet.Command, data giot_packet.Data) {
	client.queueCommand <- [2]interface{}{cmd, data}
}

func (client *ClientHandler) close() {
	client.state = CLIENT_STATE_CLOSE
}
