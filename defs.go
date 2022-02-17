package giotgo

import (
	"bytes"

	"github.com/ghuvrons/g-IoT-Go/giot_packet"
)

type ClientState uint8

const (
	CLIENT_STATE_CLOSE ClientState = iota
	CLIENT_STATE_CONNECTING
	CLIENT_STATE_CONNECT
	CLIENT_STATE_READIKNG_PAYLOAD
)

type CommandHandler func(client *ClientHandler, data giot_packet.Data) *bytes.Buffer
