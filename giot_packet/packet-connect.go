package giot_packet

import "bytes"

type PacketConnect struct {
	properties [MAX_PROP_CONNECT]Property
}

func PacketConnectDecode(pack *Packet) *PacketConnect {
	packConn := &PacketConnect{}
	Properties(packConn.properties[:]).Decode(bytes.NewReader(pack.Payload.Bytes()))
	return packConn
}

func (packConn *PacketConnect) Validate(handler func(username string, password string) bool) error {
	var username string = ""
	var password string = ""

	for _, prop := range packConn.properties {
		if prop.id == PROP_AUTH_USER {
			switch prop.data.(type) {
			case []byte:
				username = string(prop.data.([]byte))

			case *bytes.Buffer:
				username = string(prop.data.(*bytes.Buffer).Bytes())
			}

		} else if prop.id == PROP_AUTH_PASS {
			switch prop.data.(type) {
			case []byte:
				password = string(prop.data.([]byte))

			case *bytes.Buffer:
				password = string(prop.data.(*bytes.Buffer).Bytes())
			}

		} else if prop.id == 0 {
			break
		}
	}

	if handler != nil && !handler(username, password) {
		return &errAuthInvalid{}
	}

	return nil
}
