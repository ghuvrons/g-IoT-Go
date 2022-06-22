package giot_packet

type ptrAddr uint64

const PROTOCOL_VERSION = 1
const PACKET_FLAG_START = 0xC7

var PACKET_FLAG_CLOSE = [2]uint8{0x7C, 0xEA}

const MAX_PROP_CONNECT = 3
const MAX_PROP_CONNACK = 1

type PacketType byte

const (
	PACKET_TYPE_CONNECT  = 0xBA
	PACKET_TYPE_CONNACK  = 0xED
	PACKET_TYPE_DATA     = 0x01
	PACKET_TYPE_COMMAND  = 0x02
	PACKET_TYPE_RESPONSE = 0x03
)

type dataType byte

const (
	DT_BYTE   = 0x01
	DT_DINT   = 0x02
	DT_UTF8   = 0x03
	DT_BINARY = 0x04
	DT_BUFFER = 0x00
)

const (
	PROP_PROTOCOL_V     = 0x01
	PROP_AUTH_USER      = 0x10
	PROP_AUTH_PASS      = 0x11
	PROP_MAX_PACKET_LEN = 0x12
)

type Command uint16

type RespStatus byte

const (
	RESP_OK            = 0x00
	RESP_UNKNOWN_ERROR = 0xFF
)
