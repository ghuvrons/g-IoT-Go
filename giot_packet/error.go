package giot_packet

type giotError struct {
	code byte
}

func (err *giotError) Error() string {
	if byte(len(giotErrorStr)) > err.code {
		return giotErrorStr[err.code]
	}
	return "Unknown Error"
}

const (
	giot_ERR_INVALID_DATA byte = iota
	giot_ERR_INVALID_PACKET
	giot_ERR_INVALID_AUTH
	giot_ERR_BUFFER_NO_SPACE
)

var giotErrorStr = []string{
	"Invalid Data",
	"Invalid Packet",
	"Username or Password is wrong",
	"No Buffer's Space",
}
var errInvalidData = &giotError{code: giot_ERR_INVALID_DATA}
var errInvalidPacket = &giotError{code: giot_ERR_INVALID_PACKET}
var errInvalidAuth = &giotError{code: giot_ERR_INVALID_AUTH}
var errBufferNoSpace = &giotError{code: giot_ERR_BUFFER_NO_SPACE}
