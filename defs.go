package giotgo

type ClientState uint8

const (
	CLIENT_STATE_CLOSE ClientState = iota
	CLIENT_STATE_CONNECTING
	CLIENT_STATE_CONNECT
	CLIENT_STATE_READIKNG_PAYLOAD
)
